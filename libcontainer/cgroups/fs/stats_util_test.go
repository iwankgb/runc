// +build linux

package fs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/opencontainers/runc/libcontainer/cgroups"

	"github.com/sirupsen/logrus"
)

func blkioStatEntryEquals(expected, actual []cgroups.BlkioStatEntry) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("blkioStatEntries length do not match")
	}
	for i, expValue := range expected {
		actValue := actual[i]
		if expValue != actValue {
			return fmt.Errorf("Expected blkio stat entry %v but found %v", expValue, actValue)
		}
	}
	return nil
}

func expectBlkioStatsEquals(t *testing.T, expected, actual cgroups.BlkioStats) {
	if err := blkioStatEntryEquals(expected.IoServiceBytesRecursive, actual.IoServiceBytesRecursive); err != nil {
		logrus.Printf("blkio IoServiceBytesRecursive do not match - %s\n", err)
		t.Fail()
	}

	if err := blkioStatEntryEquals(expected.IoServicedRecursive, actual.IoServicedRecursive); err != nil {
		logrus.Printf("blkio IoServicedRecursive do not match - %s\n", err)
		t.Fail()
	}

	if err := blkioStatEntryEquals(expected.IoQueuedRecursive, actual.IoQueuedRecursive); err != nil {
		logrus.Printf("blkio IoQueuedRecursive do not match - %s\n", err)
		t.Fail()
	}

	if err := blkioStatEntryEquals(expected.SectorsRecursive, actual.SectorsRecursive); err != nil {
		logrus.Printf("blkio SectorsRecursive do not match - %s\n", err)
		t.Fail()
	}

	if err := blkioStatEntryEquals(expected.IoServiceTimeRecursive, actual.IoServiceTimeRecursive); err != nil {
		logrus.Printf("blkio IoServiceTimeRecursive do not match - %s\n", err)
		t.Fail()
	}

	if err := blkioStatEntryEquals(expected.IoWaitTimeRecursive, actual.IoWaitTimeRecursive); err != nil {
		logrus.Printf("blkio IoWaitTimeRecursive do not match - %s\n", err)
		t.Fail()
	}

	if err := blkioStatEntryEquals(expected.IoMergedRecursive, actual.IoMergedRecursive); err != nil {
		logrus.Printf("blkio IoMergedRecursive do not match - %v vs %v\n", expected.IoMergedRecursive, actual.IoMergedRecursive)
		t.Fail()
	}

	if err := blkioStatEntryEquals(expected.IoTimeRecursive, actual.IoTimeRecursive); err != nil {
		logrus.Printf("blkio IoTimeRecursive do not match - %s\n", err)
		t.Fail()
	}
}

func expectThrottlingDataEquals(t *testing.T, expected, actual cgroups.ThrottlingData) {
	if expected != actual {
		logrus.Printf("Expected throttling data %v but found %v\n", expected, actual)
		t.Fail()
	}
}

func expectHugetlbStatEquals(t *testing.T, expected, actual cgroups.HugetlbStats) {
	if expected != actual {
		logrus.Printf("Expected hugetlb stats %v but found %v\n", expected, actual)
		t.Fail()
	}
}

func expectMemoryStatEquals(t *testing.T, expected, actual cgroups.MemoryStats) {
	expectMemoryDataEquals(t, expected.Usage, actual.Usage)
	expectMemoryDataEquals(t, expected.SwapUsage, actual.SwapUsage)
	expectMemoryDataEquals(t, expected.KernelUsage, actual.KernelUsage)
	expectPageUsageByNUMAEquals(t, expected.PageUsageByNUMA, actual.PageUsageByNUMA)

	if expected.UseHierarchy != actual.UseHierarchy {
		logrus.Printf("Expected memory use hierarchy %v, but found %v\n", expected.UseHierarchy, actual.UseHierarchy)
		t.Fail()
	}

	for key, expValue := range expected.Stats {
		actValue, ok := actual.Stats[key]
		if !ok {
			logrus.Printf("Expected memory stat key %s not found\n", key)
			t.Fail()
		}
		if expValue != actValue {
			logrus.Printf("Expected memory stat value %d but found %d\n", expValue, actValue)
			t.Fail()
		}
	}
}

func expectMemoryDataEquals(t *testing.T, expected, actual cgroups.MemoryData) {
	if expected.Usage != actual.Usage {
		logrus.Printf("Expected memory usage %d but found %d\n", expected.Usage, actual.Usage)
		t.Fail()
	}
	if expected.MaxUsage != actual.MaxUsage {
		logrus.Printf("Expected memory max usage %d but found %d\n", expected.MaxUsage, actual.MaxUsage)
		t.Fail()
	}
	if expected.Failcnt != actual.Failcnt {
		logrus.Printf("Expected memory failcnt %d but found %d\n", expected.Failcnt, actual.Failcnt)
		t.Fail()
	}
	if expected.Limit != actual.Limit {
		logrus.Printf("Expected memory limit %d but found %d\n", expected.Limit, actual.Limit)
		t.Fail()
	}
}

func expectPageUsageByNUMAEquals(t *testing.T, expected, actual cgroups.PageUsageByNUMAWrapped) {
	if !reflect.DeepEqual(expected.Total, actual.Total) {
		logrus.Printf("Expected total page usage by NUMA %#v but found %#v", expected.Total, actual.Total)
		t.Fail()
	}
	if !reflect.DeepEqual(expected.File, actual.File) {
		logrus.Printf("Expected file page usage by NUMA %#v but found %#v", expected.File, actual.File)
		t.Fail()
	}
	if !reflect.DeepEqual(expected.Anon, actual.Anon) {
		logrus.Printf("Expected anon page usage by NUMA %#v but found %#v", expected.Anon, actual.Anon)
		t.Fail()
	}
	if !reflect.DeepEqual(expected.Unevictable, actual.Unevictable) {
		logrus.Printf("Expected unevictable page usage by NUMA %#v but found %#v", expected.Unevictable, actual.Unevictable)
		t.Fail()
	}
	if !reflect.DeepEqual(expected.Hierarchical.Total, actual.Hierarchical.Total) {
		logrus.Printf("Expected hierarchical total page usage by NUMA %#v but found %#v", expected.Hierarchical.Total, actual.Hierarchical.Total)
		t.Fail()
	}
	if !reflect.DeepEqual(expected.Hierarchical.File, actual.Hierarchical.File) {
		logrus.Printf("Expected hierarchical file page usage by NUMA %#v but found %#v", expected.Hierarchical.File, actual.Hierarchical.File)
		t.Fail()
	}
	if !reflect.DeepEqual(expected.Hierarchical.Anon, actual.Hierarchical.Anon) {
		logrus.Printf("Expected hierarchical anon page usage by NUMA %#v but found %#v", expected.Hierarchical.Anon, actual.Hierarchical.Anon)
		t.Fail()
	}
	if !reflect.DeepEqual(expected.Hierarchical.Unevictable, actual.Hierarchical.Unevictable) {
		logrus.Printf("Expected hierarchical total page usage by NUMA %#v but found %#v", expected.Hierarchical.Unevictable, actual.Hierarchical.Unevictable)
		t.Fail()
	}
}
