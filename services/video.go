package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/dato7898/grpc-tube/pb"
	"github.com/dato7898/grpc-tube/util"
	"github.com/lithammer/shortuuid/v3"
)

func (s *Server) UploadVideo(stream pb.Video_UploadVideoServer) error {
	var uf *os.File
	defer func() {
		if uf != nil {
			defer os.Remove(uf.Name())
		}
	}()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {

			tf, err := os.CreateTemp(
				"uploads",
				"grpc-tube-transcode-*.mp4",
			)
			if err != nil {
				return fmt.Errorf("error creating temporary file for transcoding: %w", err)
			}

			vf, err := securejoin.SecureJoin(
				"videos",
				fmt.Sprintf("%s.mp4", shortuuid.New()),
			)
			if err != nil {
				return fmt.Errorf("error creating file name in target library: %w", err)
			}

			// If the (sanitized) original filename collides with an existing file,
			// we try to add a shortuuid() to it until we find one that doesn't exist.
			for _, err := os.Stat(vf); !os.IsNotExist(err); _, err = os.Stat(vf) {
				if err != nil {
					return err
				}
				vf, err = securejoin.SecureJoin(
					"videos",
					fmt.Sprintf("%s_%s.mp4", filenameWithoutExtension(vf), shortuuid.New()),
				)
				if err != nil {
					return fmt.Errorf("error creating file name in target library: %w", err)
				}
			}

			thumbFn1 := fmt.Sprintf("%s.jpg", strings.TrimSuffix(tf.Name(), filepath.Ext(tf.Name())))
			thumbFn2 := fmt.Sprintf("%s.jpg", strings.TrimSuffix(vf, filepath.Ext(vf)))

			if err := util.RunCmd(
				300,
				"ffmpeg",
				"-y",
				"-i", uf.Name(),
				"-vcodec", "h264",
				"-acodec", "aac",
				"-strict", "-2",
				"-loglevel", "quiet",
				"-metadata", fmt.Sprintf("title=%s", "title"),
				"-metadata", fmt.Sprintf("comment=%s", "description"),
				tf.Name(),
			); err != nil {
				return fmt.Errorf("error transcoding video: %w", err)
			}

			if err := util.RunCmd(
				60,
				"ffmpeg",
				"-i", uf.Name(),
				"-y",
				"-vf", "thumbnail",
				"-t", fmt.Sprint(3),
				"-vframes", "1",
				"-strict", "-2",
				"-loglevel", "quiet",
				thumbFn1,
			); err != nil {
				return fmt.Errorf("error generating thumbnail: %w", err)
			}

			if err := os.Rename(thumbFn1, thumbFn2); err != nil {
				return fmt.Errorf("error renaming generated thumbnail: %w", err)
			}

			if err := os.Rename(tf.Name(), vf); err != nil {
				return fmt.Errorf("error renaming transcoded video: %w", err)
			}

			sizes := map[string]string{
				"hd720": "720p",
				"hd480": "480p",
				"nhd":   "360p",
				"film":  "240p",
			}
			for size, suffix := range sizes {
				sf := fmt.Sprintf(
					"%s#%s.mp4",
					strings.TrimSuffix(vf, filepath.Ext(vf)),
					suffix,
				)

				if err := util.RunCmd(
					300,
					"ffmpeg",
					"-y",
					"-i", vf,
					"-s", size,
					"-c:v", "libx264",
					"-c:a", "aac",
					"-crf", "18",
					"-strict", "-2",
					"-loglevel", "quiet",
					"-metadata", fmt.Sprintf("title=%s", "title"),
					"-metadata", fmt.Sprintf("comment=%s", "description"),
					sf,
				); err != nil {
					return fmt.Errorf("error transcoding video: %w", err)
				}
			}

			return stream.SendAndClose(&pb.UploadState{Success: true, Message: "File uploaded successfully"})
		}
		if err != nil {
			return err
		}

		// Open the file only on receiving the first chunk
		if uf == nil {
			uf, err = os.CreateTemp(
				"uploads",
				fmt.Sprintf("grpc-tube-upload-*%s", filepath.Ext(chunk.Filename)),
			)
			if err != nil {
				return fmt.Errorf("error creating temporary file for uploading: %w", err)
			}
		}

		// Write the current chunk to the file
		if _, err := uf.Write(chunk.Content); err != nil {
			return fmt.Errorf("failed to write chunk to file: %w", err)
		}
	}
}

func filenameWithoutExtension(path string) (stem string) {
	basename := filepath.Base(path)
	return basename[0 : len(basename)-len(filepath.Ext(basename))]
}
