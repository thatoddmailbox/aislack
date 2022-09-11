package slack

import "io"

func (c *Client) UploadFile(title string, comment string, channelID ChannelID, file io.Reader) error {
	return c.requestMultipartPOST("files.upload", map[string]string{
		"title":           title,
		"initial_comment": comment,
		"channels":        string(channelID),
	}, file, nil)
}
