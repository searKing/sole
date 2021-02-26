// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build tools

package sole

// example to get latest .git-subtree.sh
////go:generate bash -c "curl -s -L -o .git-subtree.sh https://raw.githubusercontent.com/searKing/sole/master/.git-subtree.sh"
////go:generate bash -c "chmod a+x .git-subtree.sh"
//go:generate bash -c "./.git-subtree.sh purge --remove_files"
//go:generate bash -c "./.git-subtree.sh tidy -u"
