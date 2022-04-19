package main

import (
  "fmt"
  "os"
  "github.com/akamensky/argparse"
  "bufio"
  "regexp"
  "strings"
  "strconv"
)

// splits a hex string to a int array
// e.G. #ff0000 > [255, 0, 0]
func hex_str_to_int_arr(hex_str string) [3]uint8 {

  hex_str_arr := strings.SplitAfter(hex_str, "")

  var red_hex string
  var green_hex string
  var blue_hex string
  var short bool = false

  if len(hex_str) == 7 {
    red_hex = strings.Join(hex_str_arr[1:2], "")
    green_hex = strings.Join(hex_str_arr[3:4], "")
    blue_hex = strings.Join(hex_str_arr[5:6], "")
  } else if len(hex_str) == 4 {
    red_hex = hex_str_arr[1]
    green_hex = hex_str_arr[2]
    blue_hex = hex_str_arr[3]
    short = true
  } else {
    panic(hex_str + " is no valid hex color (3 or 6 digits are needed)")
  }

  red_dec, red_err := strconv.ParseUint(red_hex, 16, 8)
  green_dec, green_err := strconv.ParseUint(green_hex, 16, 8)
  blue_dec, blue_err := strconv.ParseUint(blue_hex, 16, 8)

  if red_err != nil || green_err != nil || blue_err != nil {
    panic("Cannot not convert " + hex_str + " to int array. Ensure that arguments and data in the input file are valid hex")
  }

  var rgb_arr [3]uint8
  // If only 3 digits multiply by 17 to convert to 6 digit
  if short {
     rgb_arr[0] = uint8(red_dec) * 17
     rgb_arr[1] = uint8(green_dec) * 17
     rgb_arr[2] = uint8(blue_dec) * 17
  } else {
     rgb_arr[0] = uint8(red_dec)
     rgb_arr[1] = uint8(green_dec)
     rgb_arr[2] = uint8(blue_dec)
  }

  return rgb_arr
}

// Parse command line args
func parse_args() (string, [3]uint8) {

  parser := argparse.NewParser("darkify", "Generates dark theme CSS vars based on an existing collection of colorVars")
  input_file_path := parser.String("i", "inputFile", &argparse.Options{Required: true, Help: "File with CSS light var definitons"})
  dark_bg_str := parser.String("b", "new background", &argparse.Options{Default: "#222", Help: "The hexcode of the desired background"})

  // Parse args
  err := parser.Parse(os.Args)
  if err != nil {
    panic(parser.Usage(err))
  }

  dark_bg_rgb := hex_str_to_int_arr(*dark_bg_str)
  return *input_file_path, dark_bg_rgb
}

// Read inputfile and filter out any css vars that are not colors in hexformat (3 or 6 digits)
func read_input_file(input_file_path string) []string {

  inputFile, err := os.Open(input_file_path)
  if err != nil {
    fmt.Println("Failed to open input file")
    os.Exit(1)
  }

  scanner := bufio.NewScanner(inputFile)
  scanner.Split(bufio.ScanLines)
  var inputLines []string

  // CSS var names must be alphanum with - or _
  css_var_re, _ := regexp.Compile(`--[0-9a-zA-Z\-_]+:\s*#[0-9a-fA-F]{3}([0-9a-fA-F]{3})?\s*;`)
  for scanner.Scan() {
    var inputLine = scanner.Text()
    if css_var_re.Match([]byte(inputLine)) {
      inputLines = append(inputLines, scanner.Text())
    }
  }
  inputFile.Close()

  return inputLines
}

// Generate a map from an array containing css var definitions
func extract_css_vars(inputLines []string) map[string]string {

  var css_vars = make(map[string]string)
  for _, inputLine := range inputLines {
    lineTrimmed := strings.Trim(inputLine, " -;")
    fields := strings.Split(lineTrimmed, ":")
    css_vars[fields[0]] = strings.Trim(fields[1], " ")
  }
  return css_vars
}

func main() {

  /* default bg color, will be the replacement for white
     based on this all other colors will be calculated */
  input_file_path, dark_bg_rgb := parse_args()
  fmt.Println(dark_bg_rgb)
  inputLines := read_input_file(input_file_path)

  css_vars := extract_css_vars(inputLines)


  fmt.Println(css_vars)
}
