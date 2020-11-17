package graphics

import (
	"image"
	"sync"
)

type ImageCache struct {
	GetIfInvalid bool
	GetImage     func(string) image.Image

	key        string
	keyMutex   sync.Mutex
	image      image.Image
	imageMutex sync.Mutex
}

func (imageCache *ImageCache) Retrieve(key string) image.Image {
	imageCache.keyMutex.Lock()
	defer imageCache.keyMutex.Unlock()

	if imageCache.key != key {
		imageCache.key = key
		go imageCache.doGetImage(key)

		if !imageCache.GetIfInvalid {
			return nil
		}
	}

	return imageCache.image
}

func (imageCache *ImageCache) doGetImage(key string) {
	imageCache.imageMutex.Lock()
	defer imageCache.imageMutex.Unlock()

	imageCache.image = nil
	imageCache.image = imageCache.GetImage(key)
}
