//
// HELPERS
//


// deg2rad return radians of alpha
//
function deg2rad(alpha) {
    return (alpha/180.0) * Math.PI
}

// Clamp returns val clipped with min and max
//
function clamp(min, max, val) {
    if (val < min) {
        val = min
    }
    if (val > max) {
        val = max
    }
    return val
}

// Polar returns x,y coordinates of
// vector with raduis and alpha
//
function polar(radius, alpha) {
    return {
        x: radius * Math.cos(deg2rad(alpha)),
        y: radius * Math.sin(deg2rad(alpha))
    }
}

// Add adds to points (x,y)
//
function add(a, b) {
    a.x = a.x + b.x
    a.y = a.y + b.y
    return a
}

// Range create a list [0,1,2,...,n]
//
function range(n) {
    function idx (_,i) {
        return i
    }
    return Array
        .apply(null, Array(n))
        .map(idx)
}

// Flatten flattens a list of list into list
//
function flatten(xxs) {
    res = []
    xxs.forEach(function(xs) {
        xs.forEach(function(x){
            res.push(x)
        })
    });
    return res
}

// Bezier 3. Grades - start,end point and 2 control points
//
// 0 <= t <= 1
//
function bezierPtAt(t,p0,p1,p2,p3) {

    if (0.0 <= t && t <= 1.0) {

        /*        
        Bernstein polynoms
        ------------------
        B03 =       (1-t)*(1-t)*(1-t)
        B13 = 3  *t*(1-t)*(1-t)
        B23 = 3*t*t*(1-t)
        B33 = t*t*t
        */

        t1  = (1-t)
        t2  = t1*t1
        t3  = t2*t1
        tt3 = 3*t
        B03 =       (t3)
        B13 = tt3  *(t2)
        B23 = tt3*t*(t1)
        B33 = t*t*t

        xt = p0.x*B03 + p1.x*B13 + p2.x*B23 + p3.x*B33
        yt = p0.y*B03 + p1.y*B13 + p2.y*B23 + p3.y*B33

        return {
            x:xt, 
            y:yt
        }   
    }

    return undefined
}

// interpY interpolate y of x with points p1,p2
//
function interpY(x, p1, p2) {
    if (p1.x <= x && x <= p2.x) {
        y = x * ((p2.y - p1.y) / (p2.x - p1.x))
        return y
    }
    return undefined
}

// findIndex find the index i in a list 
// of points where pts[i].x <= x <= pts[i+1].x
//
function findIndex(x, pts) {

    function findIndexRec(x, i, pts) {
        if (i < 0 || i >= pts.length) {
            return undefined
        }

        p1 = pts[i]
        p2 = pts[i+1]
    
        if (p1.x <= x && x <= p2.x) {
            return i
        } else {
            return findIndexRec(x, i+1, pts)
        }
    }

    return findIndexRec(x, 0, pts)
}


// interpYs interpolate all x in list xs
// by using the points pts
// xs must with sorted ascending
//
function interpYs(xs, pts) {
    
    lastI = 0
    
    function calcY(x){
        i = findIndex(x, pts)
        y = interpY(x, pts[i], pts[i+1])
        lastI = i
        return y
    }

    return xs.map(calcY)
}