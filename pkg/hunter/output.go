package hunter

func PurgeDomainResult(hunterResponse ...HunterResponse) (hunterDomainResult []string) {
	if len(hunterResponse) != 0 {
		for _, response := range hunterResponse {
			for _, v := range response.Data.Arr {
				hunterDomainResult = append(hunterDomainResult, v.Url)
			}
		}
	}
	return hunterDomainResult
}
