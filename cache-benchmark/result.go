package main

import "time"

type statistic struct {
	count int
	time  time.Duration
}

/*
  getCount,setCount,missCount分别记录get，set，以及get没有找到对应
 结果的次数;
  statBuckets的下标表示这个bucket中操作花费时间的上限，statisti记录操
 作的时间与花费的总时间
*/
type result struct {
	getCount    int
	missCount   int
	setCount    int
	statBuckets []statistic
}

func (r *result) addStatistic(bucket int, stat statistic) {
	if bucket > len(r.statBuckets)-1 {
		newStatBuckets := make([]statistic, bucket+1)
		copy(newStatBuckets, r.statBuckets)
		r.statBuckets = newStatBuckets
	}
	s := r.statBuckets[bucket]
	s.count += stat.count
	s.time += stat.time
	r.statBuckets[bucket] = s
}

func (r *result) addDuration(d time.Duration, typ string) {
	bucket := int(d / time.Millisecond)
	r.addStatistic(bucket, statistic{1, d})
	if typ == "get" {
		r.getCount++
	} else if typ == "set" {
		r.setCount++
	} else {
		r.missCount++
	}
}

func (r *result) addResult(src *result) {
	for b, s := range src.statBuckets {
		r.addStatistic(b, s)
	}
	r.getCount += src.getCount
	r.setCount += src.setCount
	r.missCount += src.missCount
}
