# darkify
An algorythm to generate a dark mode from existing CSS variables

usage: darkify [-h|--help] -i|--inputFile "<value>" [-b|--new-bg "<value>"]

               Generates dark theme CSS vars based on an existing collection of
               colorVars

Arguments:

  -h  --help       Print help information
  -i  --inputFile  File with CSS light var definitons
  -b  --new-bg     The hexcode of the desired background. Default: #222
  
  
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
