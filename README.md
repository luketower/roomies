roomies
-------

`roomies` is a small command line application that calculates each roommates responsibility for a collection of bills.

### Usage
There are two ways to run the program.

  1. For one month at a time:

        roomies month 06/2015 gas 45.34 electric 47.35 rent 1200 -- bob 25 susan 25 steve 25 alice 25
  2. Read in a file of multiple months: `roomies path/to/file.txt`

        // path/to/file.txt

        month 01/2015 gas 45.34 electric 47.35 rent 1200 -- bob 25 susan 25 steve 25 alice 25
        // You can add a comment by starting a line with '//'
        // Blank lines will be ignored

        month 02/2015 gas 41.39 electric 46.48 rent 1200 -- bob 25 susan 25 steve 25 alice 25
        month 03/2015 gas 43.74 electric 43.34 rent 1200 -- bob 25 susan 25 steve 25 alice 25

        // Steve paid Alice's April bills because he owed her money.
        month 04/2015 gas 39.82 electric 37.35 rent 1200 -- bob 25 susan 25 steve 50 alice 0

The argument format is:

    roomies month <mm/yyyy> [<billname> <billamount>] -- [<roomate> <percentage>]

### Issues
Currently there is no checking on percentages, meaning if your percentages add
up to more or less than 100%, the program will still calculate. A fix is coming
soon.

I wrote this for personal use and to start learning [golang](golang.org). It's been fun. Any suggestions or nitpicks are welcome.
