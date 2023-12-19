package main

import "encoding/json"

func handleMotorData(motor Motor) {
	log(DEBUG, "Received motor data:", motor)
	// Process motor data
	byteSlice, err := json.Marshal(motor)
	if err != nil {
		log(ERROR, "Error marshalling motor data:", err)
		return
	}
	msg := "motor " + string(byteSlice) + "\n"

	(*targetConnection).Write([]byte(msg))
}

func handleLightbarData(lightbar Lightbar) {
	log(DEBUG, "Received lightbar data:", lightbar)
	// Process lightbar data
	byteSlice, err := json.Marshal(lightbar)
	if err != nil {
		log(ERROR, "Error marshalling lightbar data:", err)
		return
	}
	msg := "lightbar " + string(byteSlice) + "\n"

	(*targetConnection).Write([]byte(msg))
}

func handleBuzzerData(buzzer Buzzer) {
	log(DEBUG, "Received buzzer data:", buzzer)
	// Process buzzer data
	byteSlice, err := json.Marshal(buzzer)
	if err != nil {
		log(ERROR, "Error marshalling buzzer data:", err)
		return
	}
	msg := "buzzer " + string(byteSlice) + "\n"

	(*targetConnection).Write([]byte(msg))
}

func handleLaserdata(laser Laser) {
	log(DEBUG, "Received laser data:", laser)
	// Process buzzer data
	byteSlice, err := json.Marshal(laser)
	if err != nil {
		log(ERROR, "Error marshalling laser data:", err)
		return
	}
	msg := "laser " + string(byteSlice) + "\n"

	(*targetConnection).Write([]byte(msg))
}
