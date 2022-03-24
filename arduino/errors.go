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

package arduino

import (
	"fmt"
	"strings"

	"github.com/arduino/arduino-cli/arduino/discovery"
	"github.com/arduino/arduino-cli/i18n"
	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tr = i18n.Tr

func composeErrorMsg(msg string, cause error) string {
	if cause == nil {
		return msg
	}
	return fmt.Sprintf("%v: %v", msg, cause)
}

// CommandError is an error that may be converted into a gRPC status.
type CommandError interface {
	// ToRPCStatus convertes the error into a *status.Status
	ToRPCStatus() *status.Status
}

// InvalidInstanceError is returned if the instance used in the command is not valid.
type InvalidInstanceError struct{}

func (e *InvalidInstanceError) Error() string {
	return tr("Invalid instance")
}

// ToRPCStatus converts the error into a *status.Status
func (e *InvalidInstanceError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

// InvalidFQBNError is returned when the FQBN has syntax errors
type InvalidFQBNError struct {
	Cause error
}

func (e *InvalidFQBNError) Error() string {
	return composeErrorMsg(tr("Invalid FQBN"), e.Cause)
}

// ToRPCStatus converts the error into a *status.Status
func (e *InvalidFQBNError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

func (e *InvalidFQBNError) Unwrap() error {
	return e.Cause
}

// InvalidURLError is returned when the URL has syntax errors
type InvalidURLError struct {
	Cause error
}

func (e *InvalidURLError) Error() string {
	return composeErrorMsg(tr("Invalid URL"), e.Cause)
}

// ToRPCStatus converts the error into a *status.Status
func (e *InvalidURLError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

func (e *InvalidURLError) Unwrap() error {
	return e.Cause
}

// InvalidLibraryError is returned when the library has syntax errors
type InvalidLibraryError struct {
	Cause error
}

func (e *InvalidLibraryError) Error() string {
	return composeErrorMsg(tr("Invalid library"), e.Cause)
}

// ToRPCStatus converts the error into a *status.Status
func (e *InvalidLibraryError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

func (e *InvalidLibraryError) Unwrap() error {
	return e.Cause
}

// InvalidVersionError is returned when the version has syntax errors
type InvalidVersionError struct {
	Cause error
}

func (e *InvalidVersionError) Error() string {
	return composeErrorMsg(tr("Invalid version"), e.Cause)
}

// ToRPCStatus converts the error into a *status.Status
func (e *InvalidVersionError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

func (e *InvalidVersionError) Unwrap() error {
	return e.Cause
}

// MultipleBoardsDetectedError is returned when trying to detect
// the FQBN of a board connected to a port fails because that
// are multiple possible boards detected.
type MultipleBoardsDetectedError struct {
	Port *discovery.Port
}

func (e *MultipleBoardsDetectedError) Error() string {
	return tr(
		"Please specify an FQBN. Multiple possible ports detected on port %s with protocol %s",
		e.Port.Address,
		e.Port.Protocol,
	)
}

// ToRPCStatus converts the error into a *status.Status
func (e *MultipleBoardsDetectedError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

// MissingFQBNError is returned when the FQBN is mandatory and not specified
type MissingFQBNError struct{}

func (e *MissingFQBNError) Error() string {
	return tr("Missing FQBN (Fully Qualified Board Name)")
}

// ToRPCStatus converts the error into a *status.Status
func (e *MissingFQBNError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

// UnknownFQBNError is returned when the FQBN is not found
type UnknownFQBNError struct {
	Cause error
}

func (e *UnknownFQBNError) Error() string {
	return composeErrorMsg(tr("Unknown FQBN"), e.Cause)
}

func (e *UnknownFQBNError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *UnknownFQBNError) ToRPCStatus() *status.Status {
	return status.New(codes.NotFound, e.Error())
}

// MissingPortAddressError is returned when the port protocol is mandatory and not specified
type MissingPortAddressError struct{}

func (e *MissingPortAddressError) Error() string {
	return tr("Missing port protocol")
}

// ToRPCStatus converts the error into a *status.Status
func (e *MissingPortAddressError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

// MissingPortProtocolError is returned when the port protocol is mandatory and not specified
type MissingPortProtocolError struct{}

func (e *MissingPortProtocolError) Error() string {
	return tr("Missing port protocol")
}

// ToRPCStatus converts the error into a *status.Status
func (e *MissingPortProtocolError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

// MissingPortError is returned when the port is mandatory and not specified
type MissingPortError struct{}

func (e *MissingPortError) Error() string {
	return tr("Missing port")
}

// ToRPCStatus converts the error into a *status.Status
func (e *MissingPortError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

// NoMonitorAvailableForProtocolError is returned when a monitor for the specified port protocol is not available
type NoMonitorAvailableForProtocolError struct {
	Protocol string
}

func (e *NoMonitorAvailableForProtocolError) Error() string {
	return tr("No monitor available for the port protocol %s", e.Protocol)
}

// ToRPCStatus converts the error into a *status.Status
func (e *NoMonitorAvailableForProtocolError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

// MissingProgrammerError is returned when the programmer is mandatory and not specified
type MissingProgrammerError struct{}

func (e *MissingProgrammerError) Error() string {
	return tr("Missing programmer")
}

// ToRPCStatus converts the error into a *status.Status
func (e *MissingProgrammerError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

// ProgrammerRequiredForUploadError is returned then the upload can be done only using a programmer
type ProgrammerRequiredForUploadError struct{}

func (e *ProgrammerRequiredForUploadError) Error() string {
	return tr("A programmer is required to upload")
}

// ToRPCStatus converts the error into a *status.Status
func (e *ProgrammerRequiredForUploadError) ToRPCStatus() *status.Status {
	st, _ := status.
		New(codes.InvalidArgument, e.Error()).
		WithDetails(&rpc.ProgrammerIsRequiredForUploadError{})
	return st
}

// ProgrammerNotFoundError is returned when the programmer is not found
type ProgrammerNotFoundError struct {
	Programmer string
	Cause      error
}

func (e *ProgrammerNotFoundError) Error() string {
	return composeErrorMsg(tr("Programmer '%s' not found", e.Programmer), e.Cause)
}

func (e *ProgrammerNotFoundError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *ProgrammerNotFoundError) ToRPCStatus() *status.Status {
	return status.New(codes.NotFound, e.Error())
}

// MonitorNotFoundError is returned when the pluggable monitor is not found
type MonitorNotFoundError struct {
	Monitor string
	Cause   error
}

func (e *MonitorNotFoundError) Error() string {
	return composeErrorMsg(tr("Monitor '%s' not found", e.Monitor), e.Cause)
}

func (e *MonitorNotFoundError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *MonitorNotFoundError) ToRPCStatus() *status.Status {
	return status.New(codes.NotFound, e.Error())
}

// InvalidPlatformPropertyError is returned when a property in the platform is not valid
type InvalidPlatformPropertyError struct {
	Property string
	Value    string
}

func (e *InvalidPlatformPropertyError) Error() string {
	return tr("Invalid '%[1]s' property: %[2]s", e.Property, e.Value)
}

// ToRPCStatus converts the error into a *status.Status
func (e *InvalidPlatformPropertyError) ToRPCStatus() *status.Status {
	return status.New(codes.FailedPrecondition, e.Error())
}

// MissingPlatformPropertyError is returned when a property in the platform is not found
type MissingPlatformPropertyError struct {
	Property string
}

func (e *MissingPlatformPropertyError) Error() string {
	return tr("Property '%s' is undefined", e.Property)
}

// ToRPCStatus converts the error into a *status.Status
func (e *MissingPlatformPropertyError) ToRPCStatus() *status.Status {
	return status.New(codes.FailedPrecondition, e.Error())
}

// PlatformNotFoundError is returned when a platform is not found
type PlatformNotFoundError struct {
	Platform string
	Cause    error
}

func (e *PlatformNotFoundError) Error() string {
	return composeErrorMsg(tr("Platform '%s' not found", e.Platform), e.Cause)
}

// ToRPCStatus converts the error into a *status.Status
func (e *PlatformNotFoundError) ToRPCStatus() *status.Status {
	return status.New(codes.FailedPrecondition, e.Error())
}

func (e *PlatformNotFoundError) Unwrap() error {
	return e.Cause
}

// PlatformLoadingError is returned when a platform has fatal errors that prevents loading
type PlatformLoadingError struct {
	Cause error
}

func (e *PlatformLoadingError) Error() string {
	return composeErrorMsg(tr("Error loading hardware platform"), e.Cause)
}

// ToRPCStatus converts the error into a *status.Status
func (e *PlatformLoadingError) ToRPCStatus() *status.Status {
	return status.New(codes.FailedPrecondition, e.Error())
}

func (e *PlatformLoadingError) Unwrap() error {
	return e.Cause
}

// LibraryNotFoundError is returned when a platform is not found
type LibraryNotFoundError struct {
	Library string
	Cause   error
}

func (e *LibraryNotFoundError) Error() string {
	return composeErrorMsg(tr("Library '%s' not found", e.Library), e.Cause)
}

// ToRPCStatus converts the error into a *status.Status
func (e *LibraryNotFoundError) ToRPCStatus() *status.Status {
	return status.New(codes.FailedPrecondition, e.Error())
}

func (e *LibraryNotFoundError) Unwrap() error {
	return e.Cause
}

// LibraryDependenciesResolutionFailedError is returned when an inconsistency is found in library dependencies
// or a solution cannot be found.
type LibraryDependenciesResolutionFailedError struct {
	Cause error
}

func (e *LibraryDependenciesResolutionFailedError) Error() string {
	return composeErrorMsg(tr("No valid dependencies solution found"), e.Cause)
}

// ToRPCStatus converts the error into a *status.Status
func (e *LibraryDependenciesResolutionFailedError) ToRPCStatus() *status.Status {
	return status.New(codes.FailedPrecondition, e.Error())
}

func (e *LibraryDependenciesResolutionFailedError) Unwrap() error {
	return e.Cause
}

// PlatformAlreadyAtTheLatestVersionError is returned when a platform is up to date
type PlatformAlreadyAtTheLatestVersionError struct {
	Platform string
}

func (e *PlatformAlreadyAtTheLatestVersionError) Error() string {
	return tr("Platform '%s' is already at the latest version", e.Platform)
}

// ToRPCStatus converts the error into a *status.Status
func (e *PlatformAlreadyAtTheLatestVersionError) ToRPCStatus() *status.Status {
	st, _ := status.
		New(codes.AlreadyExists, e.Error()).
		WithDetails(&rpc.AlreadyAtLatestVersionError{})
	return st
}

// MissingSketchPathError is returned when the sketch path is mandatory and not specified
type MissingSketchPathError struct{}

func (e *MissingSketchPathError) Error() string {
	return tr("Missing sketch path")
}

// ToRPCStatus converts the error into a *status.Status
func (e *MissingSketchPathError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

// CantCreateSketchError is returned when the sketch cannot be created
type CantCreateSketchError struct {
	Cause error
}

func (e *CantCreateSketchError) Error() string {
	return composeErrorMsg(tr("Can't create sketch"), e.Cause)
}

func (e *CantCreateSketchError) Unwrap() error {
	return e.Cause
}

// CantOpenSketchError is returned when the sketch is not found or cannot be opened
type CantOpenSketchError struct {
	Cause error
}

func (e *CantOpenSketchError) Error() string {
	return composeErrorMsg(tr("Can't open sketch"), e.Cause)
}

func (e *CantOpenSketchError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *CantOpenSketchError) ToRPCStatus() *status.Status {
	return status.New(codes.NotFound, e.Error())
}

// FailedInstallError is returned if an install operation fails
type FailedInstallError struct {
	Message string
	Cause   error
}

func (e *FailedInstallError) Error() string {
	return composeErrorMsg(e.Message, e.Cause)
}

func (e *FailedInstallError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *FailedInstallError) ToRPCStatus() *status.Status {
	return status.New(codes.Internal, e.Error())
}

// FailedLibraryInstallError is returned if a library install operation fails
type FailedLibraryInstallError struct {
	Cause error
}

func (e *FailedLibraryInstallError) Error() string {
	return composeErrorMsg(tr("Library install failed"), e.Cause)
}

func (e *FailedLibraryInstallError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *FailedLibraryInstallError) ToRPCStatus() *status.Status {
	return status.New(codes.Internal, e.Error())
}

// FailedUninstallError is returned if an uninstall operation fails
type FailedUninstallError struct {
	Message string
	Cause   error
}

func (e *FailedUninstallError) Error() string {
	return composeErrorMsg(e.Message, e.Cause)
}

func (e *FailedUninstallError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *FailedUninstallError) ToRPCStatus() *status.Status {
	return status.New(codes.Internal, e.Error())
}

// FailedDownloadError is returned when a network download fails
type FailedDownloadError struct {
	Message string
	Cause   error
}

func (e *FailedDownloadError) Error() string {
	return composeErrorMsg(e.Message, e.Cause)
}

func (e *FailedDownloadError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *FailedDownloadError) ToRPCStatus() *status.Status {
	return status.New(codes.Internal, e.Error())
}

// FailedUploadError is returned when the upload fails
type FailedUploadError struct {
	Message string
	Cause   error
}

func (e *FailedUploadError) Error() string {
	return composeErrorMsg(e.Message, e.Cause)
}

func (e *FailedUploadError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *FailedUploadError) ToRPCStatus() *status.Status {
	return status.New(codes.Internal, e.Error())
}

// FailedDebugError is returned when the debug fails
type FailedDebugError struct {
	Message string
	Cause   error
}

func (e *FailedDebugError) Error() string {
	return composeErrorMsg(e.Message, e.Cause)
}

func (e *FailedDebugError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *FailedDebugError) ToRPCStatus() *status.Status {
	return status.New(codes.Internal, e.Error())
}

// FailedMonitorError is returned when opening the monitor port of a board fails
type FailedMonitorError struct {
	Cause error
}

func (e *FailedMonitorError) Error() string {
	return composeErrorMsg(tr("Port monitor error"), e.Cause)
}

func (e *FailedMonitorError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *FailedMonitorError) ToRPCStatus() *status.Status {
	return status.New(codes.Internal, e.Error())
}

// CompileFailedError is returned when the compile fails
type CompileFailedError struct {
	Message string
	Cause   error
}

func (e *CompileFailedError) Error() string {
	return composeErrorMsg(e.Message, e.Cause)
}

func (e *CompileFailedError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *CompileFailedError) ToRPCStatus() *status.Status {
	return status.New(codes.Internal, e.Error())
}

// InvalidArgumentError is returned when an invalid argument is passed to the command
type InvalidArgumentError struct {
	Message string
	Cause   error
}

func (e *InvalidArgumentError) Error() string {
	return composeErrorMsg(e.Message, e.Cause)
}

func (e *InvalidArgumentError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *InvalidArgumentError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

// NotFoundError is returned when a resource is not found
type NotFoundError struct {
	Message string
	Cause   error
}

func (e *NotFoundError) Error() string {
	return composeErrorMsg(e.Message, e.Cause)
}

func (e *NotFoundError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *NotFoundError) ToRPCStatus() *status.Status {
	return status.New(codes.NotFound, e.Error())
}

// PermissionDeniedError is returned when a resource cannot be accessed or modified
type PermissionDeniedError struct {
	Message string
	Cause   error
}

func (e *PermissionDeniedError) Error() string {
	return composeErrorMsg(e.Message, e.Cause)
}

func (e *PermissionDeniedError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *PermissionDeniedError) ToRPCStatus() *status.Status {
	return status.New(codes.PermissionDenied, e.Error())
}

// UnavailableError is returned when a resource is temporarily not available
type UnavailableError struct {
	Message string
	Cause   error
}

func (e *UnavailableError) Error() string {
	return composeErrorMsg(e.Message, e.Cause)
}

func (e *UnavailableError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *UnavailableError) ToRPCStatus() *status.Status {
	return status.New(codes.Unavailable, e.Error())
}

// TempDirCreationFailedError is returned if a temp dir could not be created
type TempDirCreationFailedError struct {
	Cause error
}

func (e *TempDirCreationFailedError) Error() string {
	return composeErrorMsg(tr("Cannot create temp dir"), e.Cause)
}

func (e *TempDirCreationFailedError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *TempDirCreationFailedError) ToRPCStatus() *status.Status {
	return status.New(codes.Unavailable, e.Error())
}

// FileCreationFailedError is returned if a file could not be created
type FileCreationFailedError struct {
	Message string
	Cause   error
}

func (e *FileCreationFailedError) Error() string {
	return composeErrorMsg(e.Message, e.Cause)
}

func (e *FileCreationFailedError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *FileCreationFailedError) ToRPCStatus() *status.Status {
	return status.New(codes.Unavailable, e.Error())
}

// SignatureVerificationFailedError is returned if a signature verification fails
type SignatureVerificationFailedError struct {
	File  string
	Cause error
}

func (e *SignatureVerificationFailedError) Error() string {
	return composeErrorMsg(tr("'%s' has an invalid signature", e.File), e.Cause)
}

func (e *SignatureVerificationFailedError) Unwrap() error {
	return e.Cause
}

// ToRPCStatus converts the error into a *status.Status
func (e *SignatureVerificationFailedError) ToRPCStatus() *status.Status {
	return status.New(codes.Unavailable, e.Error())
}

// MultiplePlatformsError is returned when trying to detect
// the Platform the user is trying to interact with and
// and multiple results are found.
type MultiplePlatformsError struct {
	Platforms    []string
	UserPlatform string
}

func (e *MultiplePlatformsError) Error() string {
	return tr("Found %d platform for reference \"%s\":\n%s",
		len(e.Platforms),
		e.UserPlatform,
		strings.Join(e.Platforms, "\n"),
	)
}

// ToRPCStatus converts the error into a *status.Status
func (e *MultiplePlatformsError) ToRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}
