# FlexyCAT
A command line CAT control for FLEX radios


Welcome to Flexy CAT - A serial CAT controller for Flex 

    by the Alan Gordon Group
           v1.0

This software is free to use, copy, and distribute. No warranty and use AS-IS. 
If you do have suggestions or issue please email me (N7AKG@ARRL.net). Have fun and 73.

--------------------------------------
Install and build
--------------------------------------
Install GO on your computer (https://go.dev/)
The GO language is supported on most PC OSs and
this program does not need any UI so should work well in any environment.

Suggested compile steps:

   1 Create a directory "build" 

   2 Compile using the GO BUILD command:
      **go build -o build/flexycat.exe cmd/flexSerialCat/main.go**

   (See the build.bat for example compile)

A pre-built Windows executable is available in 'releases'

-------------------------------------------------
Usage
-------------------------------------------------

Edit the config.ini file for port and baud rate
  or preface each command with 'CONFIG{PORT=COM8,BAUD=9600}'
  Example:
     CONFIG{PORT=COM8,BAUD=9600}:CAT:ZZMD03;ZZFA00014054350;

Example config.ini file:
```
port: COM8
baud: 9600
```

Commands:

  **CAT**:ZZxx  send CAT command to radio. Can be string together using semi-colon separator

     Note use of colon and semi-colon in command string.

  Example:
     CAT:ZZMD03;ZZFA00014054350;


  **GET**   to get response from radio and save for later send back to radio

  **SET**   send stored data to radio

   Responses are saved in the file "saved.txt" by id name for later use by the SET command

   cmd = A FLEX CAT command. e.g ZZFA;  Can string commands and save under one ID.

   Note use of colon and semi-colon in command string.

Example to GET current VFO-A settings
   GET:VFOA:ZZFA;ZZFI;ZZGT;ZZMD;ZZMG;ZZRG;ZZRT;ZZXG;

Example to send saved data to radio
   SET:VFOA

----------------------------
Examples
----------------------------

Example command to get VFO data and save it to saved.txt

```
flexycat.exe GET:SSB-10:ZZFA;ZZFI;ZZGT;ZZMD;ZZMG;ZZRG;ZZRT;ZZXG
```

-----------------------------
Example canned data for quick VFO setting. (file saved.txt)

```
CW-20:ZZFA00014030000;ZZFI04;ZZGT4;ZZMD03;ZZMG049;ZZRG+00000;ZZRT0;ZZXG+00000;
CW-40:ZZFA00007030000;ZZFI04;ZZGT4;ZZMD03;ZZMG049;ZZRG+00000;ZZRT0;ZZXG+00000;
SSB-15:ZZFA00021300000;ZZFI04;ZZGT3;ZZMD01;ZZMG049;ZZRG+00000;ZZRT0;ZZXG+00000;
CW-15:ZZFA00021030000;ZZFI04;ZZGT3;ZZMD03;ZZMG049;ZZRG+00000;ZZRT0;ZZXG+00000;
CW-10:ZZFA00028030000;ZZFI04;ZZGT3;ZZMD03;ZZMG049;ZZRG+00000;ZZRT0;ZZXG+00000;
SSB-10:ZZFA00028500000;ZZFI04;ZZGT3;ZZMD01;ZZMG049;ZZRG+00000;ZZRT0;ZZXG+00000;
SSB-40:ZZFA00007200000;ZZFI04;ZZGT3;ZZMD01;ZZMG049;ZZRG+00000;ZZRT0;ZZXG+00000;
SSB-20:ZZFA00014250000;ZZFI04;ZZGT3;ZZMD01;ZZMG049;ZZRG+00000;ZZRT0;ZZXG+00000;

```

--------------------------------------
Example use of saved data

```
flexycat.exe SET:SSB-10

or 

flexycat.exe CONFIG{PORT=COM8,BAUD=9600}:SET:SSB-10
```