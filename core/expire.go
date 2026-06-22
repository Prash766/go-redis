package core

import (
	"log"
	"math/rand"
	"time"
)

type SampleEntry struct {
	Key string
	Obj *RedisObj
}

func randomSample(count int) []SampleEntry {
	reservoir := make([]SampleEntry, 0, count)

	i := 0
	for k, obj := range store {
		entry := SampleEntry{
			Key: k,
			Obj: obj,
		}

		if i < count {
			reservoir = append(reservoir, entry)
		} else {
			j := rand.Intn(i + 1)
			if j < count {
				reservoir[j] = entry
			}
		}
		i++
	}

	return reservoir
}

func expireSampleKeys(count int) float64 {
	sample := randomSample(count)

	if len(sample) == 0 {
		return 0
	}

	expired := 0
	now := time.Now().UnixMilli()

	for _, entry := range sample {
		if entry.Obj.ExpiredAt == -1 {
			continue
		}
		if entry.Obj.ExpiredAt <= now {
			expired++
			delete(store, entry.Key)
		}
	}

	return float64(expired) / float64(len(sample))
}

func DeleteExpiredKeys(sampleCount int) {

	for {
		expiredRatio := expireSampleKeys(sampleCount)
		if expiredRatio < 0.25 {
			break
		}
	}
	log.Println("Expired keys deleted. Current store size:", len(store))
}
