package drop

import (
    "fmt"
    "os"
    "os/user"
    "strconv"
    "syscall"
)

// IsSuperuser returns true iff the current user is root
func IsSuperuser() bool {
    return os.Getuid() == 0
}

// UserLookup returns UID, GID, and Supplementary Group IDs
func UserLookup(username string) (uid int, gid int, gids []int, err error) {
    var i int
    var groups []string

    usr, err := user.Lookup(username)
    if err != nil { goto fail }

    i, err = strconv.Atoi(usr.Uid)
    if (err != nil) || (i == -1) { goto fail }
    uid = i

    i, err = strconv.Atoi(usr.Gid)
    if (err != nil) || (i == -1) { goto fail }
    gid = i

    groups, err = usr.GroupIds()
    if err != nil { goto fail }

    gids = make([]int, 0, len(groups))
    for _, g := range groups {
        i, err = strconv.Atoi(g)
        if (err != nil) || (i == -1) { goto fail }
        gids = append(gids, i)
    }

    return

    fail:
        return 0, 0, nil, fmt.Errorf("user lookup failure: %v", err)
}

// rootSetMode
func rootSetMode(path string, uid int, gid int, mode os.FileMode) (err error) {
    if !IsSuperuser() { err = fmt.Errorf("not superuser"); goto Err }

    // restrict all permissions first before changing owner so that there is
    // no race condition of enhanced permissions for anyone.

    if int(mode.Perm()) != 0 {
        err = syscall.Chmod(path, 0000)
        if err != nil { goto Err }
    }

    if (uid != -1) || (gid != -1) {
        err = syscall.Chown(path, uid, gid)
        if err != nil { goto Err }
    }

    if mode.Perm() != 0 {
        err = syscall.Chmod(path, uint32(mode.Perm()))
        if err != nil { goto Err }
    }

    return nil

    Err:
        return fmt.Errorf("error setting mode of %q: %v", path, err)
}
