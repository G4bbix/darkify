package main

// import "fmt"

func calc_rgb_sum(rgb [3]uint8) uint16 {
	var sum uint16
	for _, element := range rgb {
		sum += uint16(element)
	}
	return sum
}

func get_amount_complete_segments(rgb [3]uint8, limit uint8) uint8 {
	var complete_segments uint8 = 0
	for _, element := range rgb {
		if element == limit {
			complete_segments++
		}
	}
	return complete_segments
}

func Darkify_linear(sum_dark_rgb uint16, input_rgb [3]uint8) [3]uint8 {
	var limit uint8 = 0
	if sum_dark_rgb > calc_rgb_sum(input_rgb) {
		limit = 255
	}
	// fmt.Printf("\nLIMIT: %d\n", limit)
	var dark_rgb [3]uint8 = input_rgb
	iter := 0
	for sum_dark_rgb != calc_rgb_sum(dark_rgb) {
		// fmt.Printf("\nITER: %d\n", iter)
		iter += 1
		var amount_to_dist int16 = int16(sum_dark_rgb) - int16(calc_rgb_sum(dark_rgb))
		// fmt.Println(amount_to_dist)
		// fmt.Println(dark_rgb)
		//fmt.Println(sum_dark_rgb)
		//fmt.Println(calc_rgb_sum(dark_rgb))
		complete_segments := get_amount_complete_segments(dark_rgb, limit)
		modification_factor := float32(1) / float32(3-complete_segments)
		// fmt.Printf("MF %f\n", modification_factor)

		val_to_dist := modification_factor * float32(amount_to_dist)
		if complete_segments == 3 || (val_to_dist < 1 && val_to_dist > -1) {
			// fmt.Println("BREAK")
			break
		}
		//fmt.Printf("modfact %f\n", modification_factor)

		for i := 0; i < 3; i++ {
			if dark_rgb[i] != limit {
				// fmt.Printf("%f * %f = val_to_dist %f\n", float32(modification_factor), float32(amount_to_dist), val_to_dist)
				new_val := val_to_dist + float32(dark_rgb[i])
				if new_val >= 255 || new_val <= 0 {
					// fmt.Println("LIMIT REACHED")
					dark_rgb[i] = limit
				} else {
					dark_rgb[i] += uint8(val_to_dist)
				}
			}
		}
		//fmt.Println("leftover:")
		//fmt.Println(int16(sum_dark_rgb) - int16(calc_rgb_sum(dark_rgb)))

		if iter == 10 {
			break
		}
	}
	return dark_rgb
}
