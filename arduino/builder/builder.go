// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package builder

import (
	"errors"
	"fmt"

	"github.com/arduino/arduino-cli/arduino/builder/compilation"
	"github.com/arduino/arduino-cli/arduino/builder/detector"
	"github.com/arduino/arduino-cli/arduino/builder/logger"
	"github.com/arduino/arduino-cli/arduino/builder/progress"
	"github.com/arduino/arduino-cli/arduino/cores"
	"github.com/arduino/arduino-cli/arduino/libraries/librariesmanager"
	"github.com/arduino/arduino-cli/arduino/sketch"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
)

// ErrSketchCannotBeLocatedInBuildPath fixdoc
var ErrSketchCannotBeLocatedInBuildPath = errors.New("sketch cannot be located in build path")

// Builder is a Sketch builder.
type Builder struct {
	sketch          *sketch.Sketch
	buildProperties *properties.Map

	buildPath          *paths.Path
	sketchBuildPath    *paths.Path
	coreBuildPath      *paths.Path
	librariesBuildPath *paths.Path

	// Parallel processes
	jobs int

	// Custom build properties defined by user (line by line as "key=value" pairs)
	customBuildProperties []string

	// core related
	coreBuildCachePath *paths.Path

	logger *logger.BuilderLogger
	clean  bool

	// Source code overrides (filename -> content map).
	// The provided source data is used instead of reading it from disk.
	// The keys of the map are paths relative to sketch folder.
	sourceOverrides map[string]string

	// Set to true to skip build and produce only Compilation Database
	onlyUpdateCompilationDatabase bool
	// Compilation Database to build/update
	compilationDatabase *compilation.Database

	// Progress of all various steps
	Progress *progress.Struct

	// Sizer results
	executableSectionsSize ExecutablesFileSections

	// C++ Parsing
	lineOffset int

	targetPlatform *cores.PlatformRelease
	actualPlatform *cores.PlatformRelease

	buildArtifacts *BuildArtifacts

	*detector.SketchLibrariesDetector
	*BuildOptionsManager
}

// BuildArtifacts contains the result of various build
type BuildArtifacts struct {
	// populated by BuildCore
	coreArchiveFilePath *paths.Path
	coreObjectsFiles    paths.PathList

	// populated by BuildLibraries
	librariesObjectFiles paths.PathList

	// populated by BuildSketch
	sketchObjectFiles paths.PathList
}

// NewBuilder creates a sketch Builder.
func NewBuilder(
	sk *sketch.Sketch,
	boardBuildProperties *properties.Map,
	buildPath *paths.Path,
	optimizeForDebug bool,
	coreBuildCachePath *paths.Path,
	jobs int,
	requestBuildProperties []string,
	hardwareDirs, builtInToolsDirs, otherLibrariesDirs paths.PathList,
	builtInLibrariesDirs *paths.Path,
	fqbn *cores.FQBN,
	clean bool,
	sourceOverrides map[string]string,
	onlyUpdateCompilationDatabase bool,
	targetPlatform, actualPlatform *cores.PlatformRelease,
	useCachedLibrariesResolution bool,
	librariesManager *librariesmanager.LibrariesManager,
	libraryDirs paths.PathList,
	logger *logger.BuilderLogger,
	progressStats *progress.Struct,
) (*Builder, error) {
	buildProperties := properties.NewMap()
	if boardBuildProperties != nil {
		buildProperties.Merge(boardBuildProperties)
	}

	if buildPath != nil {
		buildProperties.SetPath("build.path", buildPath)
	}
	if sk != nil {
		buildProperties.Set("build.project_name", sk.MainFile.Base())
		buildProperties.SetPath("build.source.path", sk.FullPath)
	}
	if optimizeForDebug {
		if debugFlags, ok := buildProperties.GetOk("compiler.optimization_flags.debug"); ok {
			buildProperties.Set("compiler.optimization_flags", debugFlags)
		}
	} else {
		if releaseFlags, ok := buildProperties.GetOk("compiler.optimization_flags.release"); ok {
			buildProperties.Set("compiler.optimization_flags", releaseFlags)
		}
	}

	// Add user provided custom build properties
	customBuildProperties, err := properties.LoadFromSlice(requestBuildProperties)
	if err != nil {
		return nil, fmt.Errorf("invalid build properties: %w", err)
	}
	buildProperties.Merge(customBuildProperties)
	customBuildPropertiesArgs := append(requestBuildProperties, "build.warn_data_percentage=75")

	sketchBuildPath, err := buildPath.Join("sketch").Abs()
	if err != nil {
		return nil, err
	}
	librariesBuildPath, err := buildPath.Join("libraries").Abs()
	if err != nil {
		return nil, err
	}
	coreBuildPath, err := buildPath.Join("core").Abs()
	if err != nil {
		return nil, err
	}

	if buildPath.Canonical().EqualsTo(sk.FullPath.Canonical()) {
		return nil, ErrSketchCannotBeLocatedInBuildPath
	}

	if progressStats == nil {
		progressStats = progress.New(nil)
	}

	libsManager, libsResolver, verboseOut, err := detector.LibrariesLoader(
		useCachedLibrariesResolution, librariesManager,
		builtInLibrariesDirs, libraryDirs, otherLibrariesDirs,
		actualPlatform, targetPlatform,
	)
	if err != nil {
		return nil, err
	}
	if logger.Verbose() {
		logger.Warn(string(verboseOut))
	}

	return &Builder{
		sketch:                        sk,
		buildProperties:               buildProperties,
		buildPath:                     buildPath,
		sketchBuildPath:               sketchBuildPath,
		coreBuildPath:                 coreBuildPath,
		librariesBuildPath:            librariesBuildPath,
		jobs:                          jobs,
		customBuildProperties:         customBuildPropertiesArgs,
		coreBuildCachePath:            coreBuildCachePath,
		logger:                        logger,
		clean:                         clean,
		sourceOverrides:               sourceOverrides,
		onlyUpdateCompilationDatabase: onlyUpdateCompilationDatabase,
		compilationDatabase:           compilation.NewDatabase(buildPath.Join("compile_commands.json")),
		Progress:                      progressStats,
		executableSectionsSize:        []ExecutableSectionSize{},
		buildArtifacts:                &BuildArtifacts{},
		targetPlatform:                targetPlatform,
		actualPlatform:                actualPlatform,
		SketchLibrariesDetector: detector.NewSketchLibrariesDetector(
			libsManager, libsResolver,
			useCachedLibrariesResolution,
			onlyUpdateCompilationDatabase,
			logger,
		),
		BuildOptionsManager: NewBuildOptionsManager(
			hardwareDirs, builtInToolsDirs, otherLibrariesDirs,
			builtInLibrariesDirs, buildPath,
			sk,
			customBuildPropertiesArgs,
			fqbn,
			clean,
			buildProperties.Get("compiler.optimization_flags"),
			buildProperties.GetPath("runtime.platform.path"),
			buildProperties.GetPath("build.core.path"), // TODO can we buildCorePath ?
			logger,
		),
	}, nil
}

// GetBuildProperties returns the build properties for running this build
func (b *Builder) GetBuildProperties() *properties.Map {
	return b.buildProperties
}

// GetBuildPath returns the build path
func (b *Builder) GetBuildPath() *paths.Path {
	return b.buildPath
}

// ExecutableSectionsSize fixdoc
func (b *Builder) ExecutableSectionsSize() ExecutablesFileSections {
	return b.executableSectionsSize
}

// Preprocess fixdoc
func (b *Builder) Preprocess() error {
	b.Progress.AddSubSteps(6)
	defer b.Progress.RemoveSubSteps()
	return b.preprocess()
}

func (b *Builder) preprocess() error {
	if err := b.buildPath.MkdirAll(); err != nil {
		return err
	}

	if err := b.BuildOptionsManager.WipeBuildPath(); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.RunRecipe("recipe.hooks.prebuild", ".pattern", false); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.prepareSketchBuildPath(); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	b.logIfVerbose(false, tr("Detecting libraries used..."))
	err := b.SketchLibrariesDetector.FindIncludes(
		b.buildPath,
		b.buildProperties.GetPath("build.core.path"),
		b.buildProperties.GetPath("build.variant.path"),
		b.sketchBuildPath,
		b.sketch,
		b.librariesBuildPath,
		b.buildProperties,
		b.targetPlatform.Platform.Architecture,
	)
	if err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	b.warnAboutArchIncompatibleLibraries(b.SketchLibrariesDetector.ImportedLibraries())
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	b.logIfVerbose(false, tr("Generating function prototypes..."))
	if err := b.preprocessSketch(b.SketchLibrariesDetector.IncludeFolders()); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	// Output arduino-preprocessed source
	preprocessedSketch, err := b.sketchBuildPath.Join(b.sketch.MainFile.Base() + ".cpp").ReadFile()
	if err != nil {
		return err
	}
	b.logger.WriteStdout(preprocessedSketch)

	return nil
}

func (b *Builder) logIfVerbose(warn bool, msg string) {
	if !b.logger.Verbose() {
		return
	}
	if warn {
		b.logger.Warn(msg)
		return
	}
	b.logger.Info(msg)
}

// Build fixdoc
func (b *Builder) Build() error {
	b.Progress.AddSubSteps(6 /** preprocess **/ + 21 /** build **/)
	defer b.Progress.RemoveSubSteps()

	if err := b.preprocess(); err != nil {
		return err
	}

	buildErr := b.build()

	b.SketchLibrariesDetector.PrintUsedAndNotUsedLibraries(buildErr != nil)
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	b.printUsedLibraries(b.SketchLibrariesDetector.ImportedLibraries())
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if buildErr != nil {
		return buildErr
	}
	if err := b.exportProjectCMake(b.SketchLibrariesDetector.ImportedLibraries(), b.SketchLibrariesDetector.IncludeFolders()); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.size(); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	return nil
}

// Build fixdoc
func (b *Builder) build() error {
	b.logIfVerbose(false, tr("Compiling sketch..."))
	if err := b.RunRecipe("recipe.hooks.sketch.prebuild", ".pattern", false); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.BuildSketch(b.SketchLibrariesDetector.IncludeFolders()); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.RunRecipe("recipe.hooks.sketch.postbuild", ".pattern", true); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	b.logIfVerbose(false, tr("Compiling libraries..."))
	if err := b.RunRecipe("recipe.hooks.libraries.prebuild", ".pattern", false); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.removeUnusedCompiledLibraries(b.SketchLibrariesDetector.ImportedLibraries()); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.buildLibraries(b.SketchLibrariesDetector.IncludeFolders(), b.SketchLibrariesDetector.ImportedLibraries()); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.RunRecipe("recipe.hooks.libraries.postbuild", ".pattern", true); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	b.logIfVerbose(false, tr("Compiling core..."))
	if err := b.RunRecipe("recipe.hooks.core.prebuild", ".pattern", false); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.buildCore(); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.RunRecipe("recipe.hooks.core.postbuild", ".pattern", true); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	b.logIfVerbose(false, tr("Linking everything together..."))
	if err := b.RunRecipe("recipe.hooks.linking.prelink", ".pattern", false); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.link(); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.RunRecipe("recipe.hooks.linking.postlink", ".pattern", true); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.RunRecipe("recipe.hooks.objcopy.preobjcopy", ".pattern", false); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.RunRecipe("recipe.objcopy.", ".pattern", true); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.RunRecipe("recipe.hooks.objcopy.postobjcopy", ".pattern", true); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.MergeSketchWithBootloader(); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if err := b.RunRecipe("recipe.hooks.postbuild", ".pattern", true); err != nil {
		return err
	}
	b.Progress.CompleteStep()
	b.Progress.PushProgress()

	if b.compilationDatabase != nil {
		b.compilationDatabase.SaveToFile()
	}
	return nil
}
