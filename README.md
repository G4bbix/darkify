# darkify
An algorythm to generate a dark mode from existing CSS variables

usage: darkify [-h|--help] -i|--inputFile "<value>" [-b|--new-bg "<value>"]
 
  Generates dark theme CSS vars based on an existing collection of colorVars

  
## Arguments:

  -h  Print help information

  -i  File with CSS light var definitons

  -b  The hexcode of the desired background. Default: #000000

  -o  Output format can be hex or rgb Default: hex

  
``` 
Example output (example.css)
# ./darkify -i example.css                                [22:51:41]
--colorWhite: #222222
--colorDefault: #4d98e3
--colorRed: #220000
--colorTeal: #002222
--colorHeading: #768223
--colorBlack: #ffffff
--colorBlack2: #ffffff
```
  
  
Work in progress
