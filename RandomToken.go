package main

import (
	"math/rand"
	"time"
)

func random_token() (string, error) {
	keys := []string{
		"wxtjgoyb6atk4sbwz", "wxjificc0pvuae55u", "wxrs8pfh5o89m3o3m", "wxu9fgjmjt9axrrxv",
		"wx9vfvzqfcb2pwap1", "wxmqg2a8ack34912f", "wx5uy351nkjhdlo9g", "wxgnfokgt4hesh752",
		"wx449ef5p0ws8yymx", "wxg36g9myl4p46oll", "wxp2txz4npxkufmxy", "wxv7vgdzyc8qxai0c",
		"wxzj99nm3y0teilb8", "wxkdd3nxppsjnl4wp", "wxnhgzbdsvljv2lcm", "wxgaladufrv2y7wat",
		"wxfvnuyao7eiolf3b", "wx7ioeuno28cbevfa", "wxwjf5c0w8u464pld", "wx2k3zvhnvka6fmhc",
		"wx6me3u8cxl045wzq", "wxkh2em0pulxhhk97", "wxl2nsb5qb1l6ar2h", "wx012wq0n4l6tfbxs",
		"wxglicrnlc974babg", "wxqwqfbhy445f2gcg", "wxnx9abvq8sb6ash9", "wxhcpow8z1pld5i6s",
		"wxiqfum5rt2ja2sb5", "wxximrwdylttj1b3i",
	}

	source := rand.NewSource(time.Now().UnixNano())

	random := rand.New(source)
	rIndex := random.Intn(len(keys))

	return keys[rIndex], nil

}
