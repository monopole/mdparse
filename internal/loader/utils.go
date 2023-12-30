package loader

import "path/filepath"

// DirBase behavior:
//
//		             path  |       dir  | base
//		-------------------+------------+-----------
//		   {empty string}  |         .  |  .
//		                .  |         .  |  .
//		               ./  |         .  |  .
//	                 /  |         /  |  /
//		            ./foo  |         .  | foo
//		           ../foo  |        ..  | foo
//		             /foo  |         /  | foo
//		   /usr/local/foo  | /usr/local | foo
func DirBase(path string) (dir, base string) {
	return filepath.Dir(path), filepath.Base(path)
}
