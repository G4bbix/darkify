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
func Parse_args() (string, Rgb, string) {
	var dark_bg string
	var input_file string
	var output_format string
	flag.StringVar(&dark_bg, "d", "#000000", "Dark background (as # followed by 6 digit hex string)")
	flag.StringVar(&input_file, "f", "", "Input file containing the css vars")
	flag.StringVar(&output_format, "o", "hex", "Output format (rgb|hex)")

	flag.Parse()

	if output_format != "rgb" && output_format != "hex" {
		fmt.Println("WARNING: Invalid value for parameter -o only rgb or hex is allowed! Setting to hex")
		output_format = "hex"
	}

	dark_bg_rgb := Hex_str_to_rgb(dark_bg)
	return input_file, dark_bg_rgb, output_format
}

func float_alpha_to_uint8(alpha float64) uint8 {
	return uint8(alpha * 255)
}

// Convert a rgb(a) string (css notation) to int array
// e.G. rgba(255, 0, 0, 0.5) > Rgb{[255, 0, 0], 127}
func Css_rgb_to_rgb(rgb_str string) Rgb {
	rgb_str_arr := strings.Split(rgb_str, ",")
	var rgb Rgb
	if len(rgb_str_arr) == 3 {
		rgb.alpha = 255
	}
	for i, segment := range rgb_str_arr {
		if i <= 2 {
			segment_val, err := strconv.ParseUint(segment, 10, 8)
			if err != nil {
				fmt.Println("ERROR: Cannot parse " + rgb_str + " to int. Ensure that arguments and data in the input file are valid hex")
				break
			} else {
				rgb.color[i] = uint8(segment_val)
			}
		} else {
			segment_val, err := strconv.ParseFloat(segment, 32)
			if err != nil {
				fmt.Printf("ERROR: Cannot parse segement %f of %s to float.\n", segment_val, rgb_str)
				break
			} else {
				if segment_val <= 1 {
					rgb.alpha = float_alpha_to_uint8(segment_val)
				} else {
					fmt.Printf("ERROR: Invalid alpha value for %s\n", rgb_str)
					break
				}
			}
		}
	}
	return rgb
}

// converts a hex string to a int array
// e.G. #ff0000 > [255, 0, 0]
func Hex_str_to_rgb(hex_str_rgb string) Rgb {

	hex_str_arr := strings.SplitAfter(hex_str_rgb, "")

	var red_hex string
	var green_hex string
	var blue_hex string
	var alpha_hex string
	var short bool = false

	if len(hex_str_rgb) == 7 || len(hex_str_arr) == 9 {
		red_hex = strings.Join(hex_str_arr[1:3], "")
		green_hex = strings.Join(hex_str_arr[3:5], "")
		blue_hex = strings.Join(hex_str_arr[5:7], "")
		if len(hex_str_arr) == 9 {
			alpha_hex = strings.Join(hex_str_arr[7:9], "")
		} else {
			alpha_hex = "ff"
		}
	} else if len(hex_str_rgb) == 4 || len(hex_str_rgb) == 5 {
		red_hex = hex_str_arr[1]
		green_hex = hex_str_arr[2]
		blue_hex = hex_str_arr[3]
		short = true
		if len(hex_str_rgb) == 5 {
			alpha_hex = hex_str_arr[4]
		} else {
			alpha_hex = "f"
		}
	} else {
		panic(hex_str_rgb + " is no valid hex color (3 or 6 digits are needed)")
	}

	red_dec, red_err := strconv.ParseUint(red_hex, 16, 8)
	green_dec, green_err := strconv.ParseUint(green_hex, 16, 8)
	blue_dec, blue_err := strconv.ParseUint(blue_hex, 16, 8)
	alpha_dec, alpha_err := strconv.ParseUint(alpha_hex, 16, 8)

	if red_err != nil || green_err != nil || blue_err != nil || alpha_err != nil {
		panic("ERROR: Cannot parse " + hex_str_rgb)
	}

	var rgb Rgb
	// If only 3 digits multiply by 17 to convert to 6 digit
	if short {
		rgb.color[0] = uint8(red_dec) * 17
		rgb.color[1] = uint8(green_dec) * 17
		rgb.color[2] = uint8(blue_dec) * 17
		rgb.alpha = uint8(alpha_dec) * 17
	} else {
		rgb.color[0] = uint8(red_dec)
		rgb.color[1] = uint8(green_dec)
		rgb.color[2] = uint8(blue_dec)
		rgb.alpha = uint8(alpha_dec)
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

// Read inputfile and filter out any css vars that are colors
func Read_input_file(input_file_path string) map[string]Rgb {

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
	css_var_hex_re, _ := regexp.Compile(`--[0-9a-zA-Z\-_]+:\s*#[0-9a-fA-F]{3}([0-9a-fA-F]{3})?([0-9a-fA-F]{2})?\s*;`)
	css_var_rgb_re, _ := regexp.Compile(`--[0-9a-zA-Z\-_]+:\s*rgb\([0-9]{1,3},\s*[0-9]{1,3},\s*[0-9]{1,3}\);`)
	css_var_rgba_re, _ := regexp.Compile(`--[0-9a-zA-Z\-_]+:\s*rgba\([0-9]{1,3},\s*[0-9]{1,3},\s*[0-9]{1,3},\s*(1|0\.[0-9]{1,2})\);`)
	for scanner.Scan() {
		var input_line = scanner.Text()
		if css_var_re.Match([]byte(input_line)) {
			if css_var_rgb_re.Match([]byte(input_line)) ||
				css_var_hex_re.Match([]byte(input_line)) ||
				css_var_rgba_re.Match([]byte(input_line)) {

				input_lines = append(input_lines, scanner.Text())
			}
		}
	}
	inputFile.Close()

	input_vars := extract_css_vars(input_lines)
	return input_vars
}

// Generate a map from an array containing css var definitions
func extract_css_vars(input_lines []string) map[string]Rgb {

	var css_vars = make(map[string]Rgb)
	replacer := strings.NewReplacer(" ", "",
		"\t", "",
		"-", "",
		";", "",
		"(", "",
		")", "",
		"rgba", "",
		"rgb", "")

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
