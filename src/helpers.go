package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Rgb struct {
	name  string
	color [3]uint8
	alpha uint8
}

// Parse command line args for css
func Parse_args() (string, Rgb, string, string) {

	var dark_bg string
	var input_file string
	var output_format string
	var mode string
	flag.StringVar(&dark_bg, "d", "#000000", "Dark background (as # followed by 6 digit hex string)")
	flag.StringVar(&input_file, "f", "", "Input file containing the css vars")
	flag.StringVar(&output_format, "o", "hex", "Output format: rgb, hex")
	flag.StringVar(&mode, "m", "auto", "Mode: auto, linear, relative (Warning: relative might cause changes in hue)")

	flag.Parse()

	if output_format != "rgb" && output_format != "hex" {
		fmt.Println("WARNING: Invalid value for parameter -o only rgb or hex is allowed! Setting to hex")
		output_format = "hex"
	}

	if mode != "auto" && mode != "linear" && mode != "relative" {
		fmt.Println("WARNING: Invalid value for parameter -m only auto, linear or relative are allowed! Setting to auto")
		mode = "auto"
	}

	dark_bg_rgb := Hex_str_to_rgb(dark_bg)
	return input_file, dark_bg_rgb, output_format, mode
}
func round_float(value float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(value*ratio) / ratio
}

func Format_css_var(name string, value Rgb, format string, indent_level int) string {
	var value_formatted string
	if format == "hex" {
		value_formatted = Int_arr_to_hex_str(value.color)
		if value.alpha != 255 {
			value_formatted = fmt.Sprintf("%s%s", value_formatted, fmt.Sprintf("%x", value.alpha))
		}
	} else if format == "rgb" {
		if value.alpha == 255 {
			value_formatted = fmt.Sprintf("rgb(%d, %d, %d)", value.color[0], value.color[1], value.color[2])
		} else {
			alpha := strconv.FormatFloat(round_float(float64(value.alpha)/255, 2), 'f', -1, 32)
			value_formatted = fmt.Sprintf("rgba(%d, %d, %d, %s)", value.color[0], value.color[1], value.color[2], alpha)
		}
	}
	return fmt.Sprintf("%s--%s: %s;\n", strings.Repeat("\t", indent_level), name, value_formatted)
}

// Get the percieved birghtness of the sum of the point between the light bg and dark bg to ITU BT.601
func Calc_midpoint_sum(light_bg_rgb [3]uint8, dark_bg_rgb [3]uint8) float32 {
	return (((float32(light_bg_rgb[0]) + float32(dark_bg_rgb[0])) * 0.299) +
		((float32(light_bg_rgb[1]) + float32(dark_bg_rgb[1])) * 0.587) +
		((float32(light_bg_rgb[2]) + float32(dark_bg_rgb[2])) * 0.114)) * 3 / 2
}

// Calculate percieved brightness according to ITU BT.601
func Calc_percieved_brightness(rgb [3]uint8) uint8 {
	return uint8(float64(rgb[0])*0.299 + float64(rgb[1])*0.587 + float64(rgb[2])*0.114)
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
 * e.g. [255, 255, 255] > #ffffff
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

func Calc_remaining_unused_space(color [3]uint8) uint16 {
	var remaining_unused_space uint16
	for _, segment := range color {
		if segment != 255 {
			remaining_unused_space += uint16(255 - segment)
		}
	}
	return remaining_unused_space
}

func Calc_dark_rgb_sum(sum_light_rgb uint16, mid_point_sum float32) int16 {
	var offset int16 = int16(sum_light_rgb) - int16(mid_point_sum*2)
	// if offset is negative, invert the algabreic sign
	if offset < 0 {
		return offset / -1
	} else {
		return offset
	}
}

func Get_amount_zeros(color [3]uint8) uint8 {
	var amount_zeros uint8
	for _, segment := range color {
		if segment == 0 {
			amount_zeros++
		}
	}
	return amount_zeros
}

func Darkify_prep(light_rgb [3]uint8, mid_point_sum float32) (uint16, [3]float64) {
	sum_light_rgb := uint16(light_rgb[0]) + uint16(light_rgb[1]) + uint16(light_rgb[2])
	sum_dark_rgb := uint16(Calc_dark_rgb_sum(sum_light_rgb, mid_point_sum))
	distribution := calc_distribution(light_rgb, sum_light_rgb)
	return sum_dark_rgb, distribution
}

func calc_distribution(rgb [3]uint8, sum uint16) [3]float64 {

	if rgb[0] == 0 && rgb[1] == 0 && rgb[2] == 0 {
		var distribution = [3]float64{0.33333334, 0.33333334, 0.33333334}
		return distribution
	}

	var distribution [3]float64
	for i := 0; i < 3; i++ {
		distribution[i] = float64(rgb[i]) / float64(sum)
	}
	return distribution
}
