Despite there's no feature parity at the moment, Arduino CLI provides many of
the features you can find in the Arduino IDE, let's see some examples.

## Create a configuration file

Arduino CLI doesn't strictly require a configuration file to work because the
command line interface provides any possible functionality. However, having one
can spare you a lot of typing when issuing a command, so let's create it
right ahead with:

```sh
$ arduino-cli config init
Config file written: /home/luca/.arduino15/arduino-cli.yaml
```

If you inspect `arduino-cli.yaml` contents, you'll find out the available
options with their respective default values.

## Create a new sketch

To create a new sketch named `MyFirstSketch` in the current directory, run
the following command:

```sh
$ arduino-cli sketch new MyFirstSketch
Sketch created in: /home/luca/MyFirstSketch
```

A sketch is a folder containing assets like source files and libraries; the
`new` command creates for you a .ino file called `MyFirstSketch.ino`
containing Arduino boilerplate code:

```sh
$ cat $HOME/MyFirstSketch/MyFirstSketch.ino
void setup() {
}

void loop() {
}
```

At this point you can use your favourite file editor or IDE to open the
file `$HOME/MyFirstSketch/MyFirstSketch.ino` and change the code like this:

```c
void setup() {
    pinMode(LED_BUILTIN, OUTPUT);
}

void loop() {
    digitalWrite(LED_BUILTIN, HIGH);
    delay(1000);
    digitalWrite(LED_BUILTIN, LOW);
    delay(1000);
}
```

## Connect the board to your PC

The first thing to do upon a fresh install is to update the local cache of
available platforms and libraries by running:

```sh
$ arduino-cli core update-index
Updating index: package_index.json downloaded
```

After connecting the board to your PCs by using the USB cable, you should be
able to check whether it's been recognized by running:

```sh
$ arduino-cli board list
Port         Type              Board Name              FQBN                 Core
/dev/ttyACM1 Serial Port (USB) Arduino/Genuino MKR1000 arduino:samd:mkr1000 arduino:samd
```

In this example, the MKR1000 board was recognized and from the output of the
command you see the platform core called `arduino:samd` is the one that needs
to be installed to make it work.

If you see an `Unknown` board listed, uploading
should still work as long as you identify the platform core and use the correct
FQBN string. When a board is not detected for whatever reason, you can list all
the supported boards and their FQBN strings by running the following:

```sh
$ arduino-cli board listall mkr
Board Name              FQBN
Arduino MKR FOX 1200    arduino:samd:mkrfox1200
Arduino MKR GSM 1400    arduino:samd:mkrgsm1400
Arduino MKR WAN 1300    arduino:samd:mkrwan1300
Arduino MKR WiFi 1010   arduino:samd:mkrwifi1010
Arduino MKRZERO         arduino:samd:mkrzero
Arduino/Genuino MKR1000 arduino:samd:mkr1000
```

## Install the core for your board

To install the ``arduino:samd`` platform core, run the following:

```sh
$ arduino-cli core install arduino:samd
Downloading tools...
arduino:arm-none-eabi-gcc@4.8.3-2014q1 downloaded
arduino:bossac@1.7.0 downloaded
arduino:openocd@0.9.0-arduino6-static downloaded
arduino:CMSIS@4.5.0 downloaded
arduino:CMSIS-Atmel@1.1.0 downloaded
arduino:arduinoOTA@1.2.0 downloaded
Downloading cores...
arduino:samd@1.6.19 downloaded
Installing tools...
Installing platforms...
Results:
arduino:samd@1.6.19 - Installed
arduino:arm-none-eabi-gcc@4.8.3-2014q1 - Installed
arduino:bossac@1.7.0 - Installed
arduino:openocd@0.9.0-arduino6-static - Installed
arduino:CMSIS@4.5.0 - Installed
arduino:CMSIS-Atmel@1.1.0 - Installed
arduino:arduinoOTA@1.2.0 - Installed
```

Now verify we have installed the core properly by running:

```sh
$ arduino-cli core list
ID              Installed       Latest  Name
arduino:samd    1.6.19          1.6.19  Arduino SAMD Boards (32-bits ARM Cortex-M0+)
```

Great! Now we are ready to compile and upload the sketch.

## Adding 3rd party cores

If your board requires 3rd party core packages to work, you can list the URLs
to additional package indexes in the Arduino CLI configuration file.

For example, to add the ESP8266 core, edit the configuration file and change the
`board_manager` settings as follows:

```yaml
board_manager:
    additional_urls:
    - https://arduino.esp8266.com/stable/package_esp8266com_index.json
```

From now on, commands supporting custom cores will automatically use the
additional URL from the configuration file:

```sh
$ arduino-cli core update-index
Updating index: package_index.json downloaded
Updating index: package_esp8266com_index.json downloaded
Updating index: package_index.json downloaded

$ arduino-cli core search esp8266
ID              Version Name
esp8266:esp8266 2.5.2   esp8266
```

Alternatively, you can pass a link to the the additional package index file with
the `--additional-urls` option, that has to be specified every time and for every
command that operates on a 3rd party platform core, for example:

```sh
$ arduino-cli  core update-index --additional-urls https://arduino.esp8266.com/stable/package_esp8266com_index.json
Updating index: package_esp8266com_index.json downloaded

$ arduino-cli core search esp8266 --additional-urls https://arduino.esp8266.com/stable/package_esp8266com_index.json
ID              Version Name
esp8266:esp8266 2.5.2   esp8266
```

## Compile and upload the sketch

To compile the sketch you run the `compile` command passing the proper FQBN
string:

```sh
$ arduino-cli compile --fqbn arduino:samd:mkr1000 MyFirstSketch
Sketch uses 9600 bytes (3%) of program storage space. Maximum is 262144 bytes.
```

To upload the sketch to your board, run the following command, this time also
providing the serial port where the board is connected:

```sh
$ arduino-cli upload -p /dev/ttyACM0 --fqbn arduino:samd:mkr1000 MyFirstSketch
No new serial port detected.
Atmel SMART device 0x10010005 found
Device       : ATSAMD21G18A
Chip ID      : 10010005
Version      : v2.0 [Arduino:XYZ] Dec 20 2016 15:36:43
Address      : 8192
Pages        : 3968
Page Size    : 64 bytes
Total Size   : 248KB
Planes       : 1
Lock Regions : 16
Locked       : none
Security     : false
Boot Flash   : true
BOD          : true
BOR          : true
Arduino      : FAST_CHIP_ERASE
Arduino      : FAST_MULTI_PAGE_WRITE
Arduino      : CAN_CHECKSUM_MEMORY_BUFFER
Erase flash
done in 0.784 seconds

Write 9856 bytes to flash (154 pages)
[==============================] 100% (154/154 pages)
done in 0.069 seconds

Verify 9856 bytes of flash with checksum.
Verify successful
done in 0.009 seconds
CPU reset.
```

## Add libraries

If you need to add more functionalities to your sketch, chances are some of the
libraries available in the Arduino ecosystem already provide what you need.
For example, if you need a debouncing strategy to better handle button inputs,
you can try searching for the `debouncer` keyword:

```sh
$ arduino-cli lib search debouncer
Name: "Debouncer"
    Author: hideakitai
    Maintainer: hideakitai
    Sentence: Debounce library for Arduino
    Paragraph: Debounce library for Arduino
    Website: https://github.com/hideakitai
    Category: Timing
    Architecture: *
    Types: Contributed
    Versions: [0.1.0]
Name: "FTDebouncer"
    Author: Ubi de Feo
    Maintainer: Ubi de Feo, Sebastian Hunkeler
    Sentence: An efficient, low footprint, fast pin debouncing library for Arduino
    Paragraph: This pin state supervisor manages debouncing of buttons and handles transitions between LOW and HIGH state, calling a function and notifying your code of which pin has been activated or deactivated.
    Website: https://github.com/ubidefeo/FTDebouncer
    Category: Uncategorized
    Architecture: *
    Types: Contributed
    Versions: [1.3.0]
Name: "SoftTimer"
    Author: Balazs Kelemen <prampec+arduino@gmail.com>
    Maintainer: Balazs Kelemen <prampec+arduino@gmail.com>
    Sentence: SoftTimer is a lightweight pseudo multitasking solution for Arduino.
    Paragraph: SoftTimer enables higher level Arduino programing, yet easy to use, and lightweight. You are often faced with the problem that you need to do multiple tasks at the same time. In SoftTimer, the programmer creates Tasks that runs periodically. This library comes with a collection of handy tools like blinker, pwm, debouncer.
    Website: https://github.com/prampec/arduino-softtimer
    Category: Timing
    Architecture: *
    Types: Contributed
    Versions: [3.0.0, 3.1.0, 3.1.1, 3.1.2, 3.1.3, 3.1.5, 3.2.0]
```

Our favourite is ``FTDebouncer``, let's install it by running:

```sh
$ arduino-cli lib install FTDebouncer
FTDebouncer depends on FTDebouncer@1.3.0
Downloading FTDebouncer@1.3.0...
FTDebouncer@1.3.0 downloaded
Installing FTDebouncer@1.3.0...
Installed FTDebouncer@1.3.0
```