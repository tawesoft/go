package start

import (
    "fmt"
    "os/user"
    "strconv"
)

// returns UID, GID, and Supplementary Group IDs
func userLookup(username string) (uid uint32, gid uint32, gids []uint32, err error) {
    var i int
    var groups []string
    
    usr, err := user.Lookup(username)
    if err != nil { goto fail }
    
    i, err = strconv.Atoi(usr.Uid)
    if (err != nil) || (i < 0) { goto fail }
    uid = uint32(i)
    
    i, err = strconv.Atoi(usr.Gid)
    if (err != nil) || (i < 0) { goto fail }
    gid = uint32(i)
    
    groups, err = usr.GroupIds()
    if err != nil { goto fail }
    
    gids = make([]uint32, 0, len(groups))
    for _, g := range groups {
        i, err = strconv.Atoi(g)
        if (err != nil) || (i < 0) { goto fail }
        gids = append(gids, uint32(i))
    }
    
    return
    
    fail:
        return 0, 0, nil, fmt.Errorf("user lookup failure: %v", err)
}
