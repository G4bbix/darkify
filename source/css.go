package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func float_alpha_to_uint8(alpha float64) uint8 {
	return uint8(alpha * 255)
}

// Read inputfile and filter out any css vars that are colors
func Read_input_file(input_file_path string) []Rgb {

	inputFile, err := os.Open(input_file_path)
	if err != nil {
		fmt.Println("Failed to open input file")
		os.Exit(1)
	}
	defer inputFile.Close()

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

				input_lines = append(input_lines, input_line)
			}
		}
	}
	inputFile.Close()
	input_vars := extract_css_vars(input_lines)
	return input_vars
}

// Generate a map from an array containing css var definitions
func extract_css_vars(input_lines []string) []Rgb {

	var css_vars []Rgb
	replacer := strings.NewReplacer(" ", "",
		"\t", "",
		"-", "",
		";", "",
		"(", "",
		")", "",
		"rgba", "",
		"rgb", "")

	var rgb Rgb
	for _, input_line := range input_lines {
		line_trimmed := replacer.Replace(input_line)
		rgb_name_and_val := strings.Split(line_trimmed, ":")
		rgb_name := rgb_name_and_val[0]
		rgb_val := rgb_name_and_val[1]
		if strings.HasPrefix(rgb_val, "#") {
			rgb = Hex_str_to_rgb(rgb_val)
		} else {
			rgb = Css_rgb_to_rgb(rgb_val)
		}
		rgb.name = rgb_name
		css_vars = append(css_vars, rgb)
	}

	return css_vars
}
