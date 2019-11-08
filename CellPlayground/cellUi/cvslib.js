//
// SHAPES
//
/*
var paper
var size

function init() {
    paper = getPaper()
    size  = getSize() 

    // Defaults
    //
    paper.pointSize = 1.0
    paper.DEBUG = true
}
*/

// Styles
function getStyle() {

    st = {strokeColor:'grey', lineWidth:1}
    sp = {strokeColor:'blue', lineWidth:1, fillColor: 'blue', pointSize: 1.2}
    
    return {
        tangentLine  : st,
        tangentPoint : sp
    }
}

function getPaper() {
    return document.getElementById("scene").getContext("2d")
}

function getSize() {
    width  = document.getElementById("scene").width
    height = document.getElementById("scene").height
    return {width:width, height:height}
}

// point draw a point with given point size
//
function point(pt) {
    return function() {
        paper.beginPath()
        paper.arc(pt.x, pt.y, paper.pointSize, 0, Math.PI *2, true)
        paper.fill()
        paper.stroke()
    }
}

//circle draws a with given radius
//
function circle(r) {
    return function() {
        paper.beginPath()
        paper.arc(0, 0, r, 0, Math.PI *2, true)
        paper.fill()
        paper.stroke()
    };
}

// arc draws an arc 
// with start angle a and end angle b
//
function arc(r, a, b) {
    return function() {
        paper.beginPath()
        paper.arc(0, 0, r, deg2rad(a), deg2rad(b), false)
        paper.stroke()
    }
}

// rect draws an rectangle 
// with width w and height h
//
function rect(w, h) {
    return function() {
        paper.fillRect  (-w/2, -h/2, w, h)
        paper.strokeRect(-w/2, -h/2, w, h)
    }
}

// line draws a line 
// from point p1 to p2
//
function line(p1,p2) {
    return function() {
        paper.beginPath()
        paper.moveTo(p1.x,p1.y)
        paper.lineTo(p2.x,p2.y)
        paper.stroke()
    }
}

// pline draws a polyline 
// with given list of points pts
//
function pline(pts) {
    return function() {
        function fn(pt,i){
            if (i === 0) {
                paper.moveTo(pt.x,pt.y)
            } else {
                paper.lineTo(pt.x,pt.y)
            }
        }
        paper.beginPath()
        pts.forEach(fn)
        paper.stroke()
    }
}

// bezCurve draws a bezier 3rd grade
// from point p1 to point p2
// with control points p1 and p2
//
function bezCurve(p1,cp1,cp2,p2) {
    fns = [
        function() {
            paper.beginPath()
            
            paper.moveTo(p1.x,p1.y)
            paper.bezierCurveTo(cp1.x,cp1.y,cp2.x,cp2.y,p2.x,p2.y)
            paper.stroke()
        }
    ]

    // Draw tangents in debug mode
    if (paper.DEBUG) {

        tangents = [
            style(getStyle().tangentLine,line(p1,cp1)), // in
            style(getStyle().tangentLine,line(cp2,p2)), // out
            style(getStyle().tangentPoint,point(cp1)),   // in
            style(getStyle().tangentPoint,point(cp2)),   // out
        ]

        fns = tangents.concat(fns)
    }

    return group(fns)
}

// quadCurve draws a quadratic curve
// from points p1 to points p2
// with given control point cp
//
function quadCurve(p1,cp,p2) {

    fn = function() {
        paper.beginPath()
        paper.moveTo(p1.x,p1.y)
        paper.quadraticCurveTo(cp.x,cp.y,p2.x,p2.y)
        paper.stroke()
    }

    fns = [fn]

    function quadCurveTangents() {
   
        tangents = [
            style(getStyle().tangentLine,line(p1,cp)), // in
            style(getStyle().tangentLine,line(cp,p2)), // out
            style(getStyle().tangentPoint,point(cp)),
        ]
        return tangents
    }

    if (paper.DEBUG) {
        tangents = quadCurveTangents()
        fns = tangents.concat(fns)  
    }

    return group(fns)
}

// text draws a text from given string str
//
function text(str) {
    return function() {
        paper.fillText  (str, 0, 0)
        paper.strokeText(str, 0, 0)
    }
}

//
// MODIFIERS
//

// group groups a list of draw functions
// and execute each one
//
function group(fns) {
    return function() {
        fns.forEach(function(fn) {
            fn()        
        })
    }
}

// transform transfrom to t.x, t.y
// and executes draw fn
//
function transform(t, fn) {
    return function() {
        paper.save()
            paper.translate(t.x,  t.y )
            paper.scale    (t.sx, t.sy)
            fn()
        paper.restore()
    }
}

// style styles with style s
// and executes draw fn
//
function style(s, fn) {
    return function() {
        
        if (s.pointSize !== undefined) {
            savedPointSize = paper.pointSize
        }

        paper.save()

            paper.font         = s.font
            paper.textAlign    = s.textAlign
            paper.textBaseline = s.textBaseline

            paper.fillStyle    = s.fillColor
            paper.strokeStyle  = s.strokeColor
            
            paper.lineWidth    = s.lineWidth

            if (s.pointSize !== undefined) {
                paper.pointSize  = s.pointSize
            }

            if (s.lineDash !== undefined) {
                paper.setLineDash(s.lineDash)
            }

            if (s.lineDashOffset !== undefined) {
                paper.lineDashOffset = s.lineDashOffset
            }
            
            fn()

        paper.restore()

        if (s.pointSize !== undefined) {
            paper.pointSize = savedPointSize
        }
    }
}
