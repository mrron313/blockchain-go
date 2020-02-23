package timezone

import (
    "time"
)

func TimeIn(t time.Time, name string) (time.Time) {
    loc, err := time.LoadLocation(name)
    if err == nil {
        t = t.In(loc)
    }

    return t
}
