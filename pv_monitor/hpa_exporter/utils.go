package main

/**************************************************/

type PodInfo struct {
	PVCNames              []string
	PVNames               []string
	CpuMilliLimit         int64
	MemoryByteLimit       int64
}

func (p *PodInfo)AppendPVCName(PVCName string) {
	p.PVCNames = append(p.PVCNames, PVCName)
}
func (p PodInfo)GetPVCNames() []string {
	return p.PVCNames
}
func (p *PodInfo)AppendPVName(PVName string) {
	p.PVNames = append(p.PVNames, PVName)
}
func (p PodInfo)GetPVNames() []string {
	return p.PVNames
}

/**************************************************/

type PVInfo struct {
	PVDiskUtilization    float32
	PVDiskIOPS           float32
	PVDiskWriteKBPS      float32
	PVDiskReadKBPS       float32
}

func (p *PVInfo) SetPVDiskUtilization(utilization float32) {
	p.PVDiskUtilization = utilization
}
func (p PVInfo) GetPVDiskUtilization() float32 {
	return p.PVDiskUtilization
}
func (p *PVInfo) SetPVDiskIOPS(iops float32) {
	p.PVDiskIOPS = iops
}
func (p PVInfo) GetPVDiskIOPS() float32  {
	return p.PVDiskIOPS
}

/****************************************************/

type StatefulSetInfo struct {
	StatefulSetName      string                /* statefulSet name        */
	PodInfos             map[string]PodInfo    /* podName --> PodInfo     */
	PVInfos              map[string]PVInfo     /* podName --> PVInfo      */
	Initialized          bool                  /* whether the obj has been initialized */
}

func getStatefulSetInfoObj(stsName string) StatefulSetInfo {
	var stsInfo StatefulSetInfo

	stsInfo.StatefulSetName = stsName
	stsInfo.PodInfos = make(map[string]PodInfo)
	stsInfo.PVInfos  = make(map[string]PVInfo)

	return stsInfo
}
func (s *StatefulSetInfo) GetStatefulSetName() string {
	return s.StatefulSetName
}
func (s *StatefulSetInfo) SetStatefulSetName(statefulSetName string) {
	s.StatefulSetName = statefulSetName
}
func (s *StatefulSetInfo) GetPodNames() []string {
	var podNames []string

	if s.PodInfos == nil {
		s.PodInfos = make(map[string]PodInfo)
	}
	for podName, _ := range s.PodInfos {
		podNames = append(podNames, podName)
	}

	return podNames
}
func (s *StatefulSetInfo) GetPodInfo(podName string) PodInfo {
	if s.PodInfos == nil {
		s.PodInfos = make(map[string]PodInfo)
		return PodInfo{}
	}

	podInfo, found := s.PodInfos[podName]
	if found == true {
		return podInfo
	}
	return PodInfo{}
}
func (s *StatefulSetInfo) GetPodInfos() map[string]PodInfo {
	return s.PodInfos
}
func (s *StatefulSetInfo) SetPodInfo(podName string, podInfo PodInfo) {
	if s.PodInfos == nil {
		s.PodInfos = make(map[string]PodInfo)
	}

	s.PodInfos[podName] = podInfo
}
func (s *StatefulSetInfo) GetPVInfo(podName string) PVInfo {
	if s.PVInfos == nil {
		s.PVInfos = make(map[string]PVInfo)
		return PVInfo{}
	}

	pvInfo, found := s.PVInfos[podName]
	if found == true {
		return pvInfo
	}
	return PVInfo{}
}
func (s *StatefulSetInfo) GetPVInfos() map[string]PVInfo {
	return s.PVInfos
}
func (s *StatefulSetInfo) SetPVInfo(podName string, pvInfo PVInfo) {
	if s.PVInfos == nil {
		s.PVInfos = make(map[string]PVInfo)
	}
	s.PVInfos[podName] = pvInfo
}

/*************************************************/

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func getAvgFloat64Utilization(utilizations []float64) float64 {
	var utilizationSum float64 = 0
	for _, utilization := range utilizations {
		utilizationSum += utilization
	}

	return utilizationSum / float64(len(utilizations))
}

func getAvgInt64Utilization(utilizations []int64) int64 {
	var utilizationSum int64 = 0
	for _, utilization := range utilizations {
		utilizationSum += utilization
	}

	return utilizationSum / int64(len(utilizations))
}

func getAboveUtilizationNumber(utilizations []float64, cmpUtilization float64) int {
	var num = 0
	for _, utilization := range utilizations {
		if utilization > cmpUtilization {
			num++
		}
	}

	return num
}

/************************************************/

