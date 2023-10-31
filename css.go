package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Parse command line args for css
func Parse_args() (string, [3]uint8, string) {
	var dark_bg string
	var input_file string
	var output_format string
	flag.StringVar(&dark_bg, "d", "#000000", "Dark background (as # followed by 6 digit hex string)")
	flag.StringVar(&input_file, "f", "", "Input file containing the css vars")
	flag.StringVar(&output_format, "o", "hex", "Output format (rgb|hex)")

	flag.Parse()

	if output_format != "rgb" && output_format != "hex" {
		panic("Invalid value for parameter -o only rgb or hex is allowed")
	}

	dark_bg_rgb := Hex_str_to_rgb(dark_bg)
	return input_file, dark_bg_rgb, output_format
}

// Convert a rgb(a) string (css notation) to int array
// e.G. rgb(255, 0, 0) > [255, 0, 0]
func Css_rgb_to_rgb(rgb_str string) [3]uint8 {
	rgb_str_arr := strings.Split(rgb_str, ",")
	var rgb [3]uint8
	for i, segment := range rgb_str_arr {
		segment_int, err := strconv.ParseUint(segment, 10, 8)
		if err != nil {
			panic("Cannot parse " + rgb_str + " to int array. Ensure that arguments and data in the input file are valid hex")
		}
		rgb[i] = uint8(segment_int)
	}
	return rgb
}

// converts a hex string to a int array
// e.G. #ff0000 > [255, 0, 0]
func Hex_str_to_rgb(hex_str_rgb string) [3]uint8 {

	hex_str_arr := strings.SplitAfter(hex_str_rgb, "")

	var red_hex string
	var green_hex string
	var blue_hex string
	var short bool = false

	if len(hex_str_rgb) == 7 {
		red_hex = strings.Join(hex_str_arr[1:3], "")
		green_hex = strings.Join(hex_str_arr[3:5], "")
		blue_hex = strings.Join(hex_str_arr[5:7], "")
	} else if len(hex_str_rgb) == 4 {
		red_hex = hex_str_arr[1]
		green_hex = hex_str_arr[2]
		blue_hex = hex_str_arr[3]
		short = true
	} else {
		panic(hex_str_rgb + " is no valid hex color (3 or 6 digits are needed)")
	}

	red_dec, red_err := strconv.ParseUint(red_hex, 16, 8)
	green_dec, green_err := strconv.ParseUint(green_hex, 16, 8)
	blue_dec, blue_err := strconv.ParseUint(blue_hex, 16, 8)

	if red_err != nil || green_err != nil || blue_err != nil {
		panic("Cannot parse " + hex_str_rgb + " to int array. Ensure that arguments and data in the input file are valid hex")
	}

	var rgb [3]uint8
	// If only 3 digits multiply by 17 to convert to 6 digit
	if short {
		rgb[0] = uint8(red_dec) * 17
		rgb[1] = uint8(green_dec) * 17
		rgb[2] = uint8(blue_dec) * 17
	} else {
		rgb[0] = uint8(red_dec)
		rgb[1] = uint8(green_dec)
		rgb[2] = uint8(blue_dec)
	}

	return rgb
}

/* calculate an hex string from an int array
 * e.g. [255, 255, 255] > #ff0000
 */
func Int_arr_to_hex_str(input [3]uint8) string {
	var return_val string = "#"
	for _, octett := range input {
		if octett < 16 {
			return_val += "0"
		}
		return_val += fmt.Sprintf("%x", octett)
	}
	return return_val
}

// Read inputfile and filter out any css vars that are not colors in hexformat (3 or 6 digits)
func Read_input_file(input_file_path string) map[string][3]uint8 {

	inputFile, err := os.Open(input_file_path)
	if err != nil {
		fmt.Println("Failed to open input file")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)
	var input_lines []string

	// CSS var names must be alphanum with - or _
	css_var_re, _ := regexp.Compile(`--[0-9a-zA-Z\-_]+:\s*.*;`)
	css_var_hex_re, _ := regexp.Compile(`--[0-9a-zA-Z\-_]+:\s*#[0-9a-fA-F]{3}([0-9a-fA-F]{3})?\s*;`)
	css_var_rgb_re, _ := regexp.Compile(`--[0-9a-zA-Z\-_]+:\s*rgba?\([0-9]{1,3},\s*[0-9]{1,3},\s*[0-9]{1,3}\);`)
	for scanner.Scan() {
		var input_line = scanner.Text()
		if css_var_re.Match([]byte(input_line)) {
			if css_var_rgb_re.Match([]byte(input_line)) || css_var_hex_re.Match([]byte(input_line)) {
				input_lines = append(input_lines, scanner.Text())
			}
		}
	}
	inputFile.Close()

	input_vars := extract_css_vars(input_lines)
	return input_vars
}

// Generate a map from an array containing css var definitions
func extract_css_vars(input_lines []string) map[string][3]uint8 {

	var css_vars = make(map[string][3]uint8)
	replacer := strings.NewReplacer(" ", "",
		"\t", "",
		"-", "",
		";", "",
		"(", "",
		")", "",
		"rgb", "",
		"rgba", "")

	for _, input_line := range input_lines {
		line_trimmed := replacer.Replace(input_line)
		rgb_name_and_val := strings.Split(line_trimmed, ":")
		rgb_name := rgb_name_and_val[0]
		rgb_val := rgb_name_and_val[1]
		if strings.HasPrefix(rgb_val, "#") {
			css_vars[rgb_name] = Hex_str_to_rgb(rgb_val)
		} else {
			css_vars[rgb_name] = Css_rgb_to_rgb(rgb_val)
		}
	}
	return css_vars
}
