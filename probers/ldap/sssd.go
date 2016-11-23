package ldap

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func Makeldapprober(testname string, sssdconf string,
	probefreq uint8) *ldapconfig {
	comment := regexp.MustCompile("#+.*")
	uri := regexp.MustCompile("ldap_uri ?= ?(.*)")
	bind := regexp.MustCompile("ldap_default_bind_dn ?= ?(.*)")
	tok := regexp.MustCompile("ldap_default_authtok ?= ?(.*)")

	conf, err := os.Open(sssdconf)
	if err != nil {
		panic(err)
	}
	defer conf.Close()
	var hostnames []string
	var val, binddn, bindpwd string
	scanner := bufio.NewScanner(conf)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		l := scanner.Text()
		cmt := comment.FindStringSubmatch(l)
		if len(cmt) > 0 {
			continue
		}
		val = findmatch(uri.FindStringSubmatch(l))
		if val != "" {
			hostnames = strings.Split(val, ", ")
			for i, host := range hostnames {
				hostnames[i] = strings.TrimPrefix(host, "ldap://")
			}
			continue
		}
		val = findmatch(bind.FindStringSubmatch(l))
		if val != "" {
			binddn = val
			continue
		}
		val = findmatch(tok.FindStringSubmatch(l))
		if val != "" {
			bindpwd = val
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return New(testname, probefreq, hostnames, binddn, bindpwd)
}

func findmatch(list []string) string {
	if len(list) > 0 {
		return list[1]
	}
	return ""
}
