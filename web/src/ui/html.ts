
import jQuery from "jquery";

const $ = jQuery;

type Attr = string | (()=>string) ;

type H = Attr | (()=>JQuery) | JQuery;

const str = (s: Attr) => {
    return typeof(s)=="function" ? s(): s;
}

const to = (s:H) => {
    return typeof(s)=="function" ? s(): s;
}

export const h = (name:string, attrs?: Record<string, Attr>, ...children: H[]): JQuery => {

    const $e = $('<' + name +'>');
    if(attrs) {
        Object.entries(attrs).forEach(v => {
            $e.attr(v[0], str(v[1]));
        });
    }
    if(children) {
        children.forEach(v => {
            $e.append( to(v) );
        });
    }
    return $e;
}
