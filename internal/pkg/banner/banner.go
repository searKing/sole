// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package banner

func Banner(name, version string) string {
	return `Thank you for using ` + name + `` + version + `!
Take security seriously and subscribe to the searKing Github Issue. Stay on top of new patches and security insights. 
>> Subscribe now: https://github.com/searKing/sole/issues <<`
}
