# Greg
(Homage to the taskmaster)

Cross Platform (written in GO) parse for Windows XML task files.  Manually created the structure from samples I had - seems to be work fine. (I am not responsible for missing output or parsing goofs- as a forensicator you need to validate your own tools used in process)

Takes directory of Tasks as an input parameter and can be output to screen -s and/or a csv with -o flag

Check releases for pre-compiled version (this wont always be updated in line with other changes) - but will do my best
<code>
-----------

Usage

 _____
|  __ \
| |  \/ _ __   ___   __ _
| | __ | '__| / _ \ / _` |
| |_\ \| |   |  __/| (_| |
 \____/|_|    \___| \__, |
                     __/ |
                    |___/

The Taskmaster - Windows Task XML parser

usage: Okonma [-h|--help] -i|--inputdir "<value>" [-o|--outputdir "<value>"]
              [-s|--show]

              Go Wrapper for Forensic Tools

Arguments:

  -h  --help       Print help information
  -i  --inputdir   Input Directory
  -o  --outputdir  (optional) param. Output Directory
  -s  --show       Dont Output to a file - show me the results onscreen
  
</code>

example

<img width="3205" alt="image" src="https://user-images.githubusercontent.com/22748755/205974950-c9cad41a-03e8-4aea-aa8f-c978002cdb6c.png">

