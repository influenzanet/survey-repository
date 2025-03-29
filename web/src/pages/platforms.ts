import { apiLoadPlatforms } from "../api";
import { h } from "../ui";
import { Page } from "./page";

export class Platform extends Page {

    constructor() {
        super('platforms');
    }

    render(): void {
        const $tb = $('#platforms-rows');
        $tb.empty();

        apiLoadPlatforms().then(platforms => {
            platforms.forEach(p => {
                const $r = h('tr');
                $r.append( h('td', undefined,  h('span', {"class":"badge bg-danger me-1"}, p.id) ) );
                $r.append( h('td', undefined, p.label));
                $r.append( h('td', undefined, p.country));
                $tb.append($r);
            });
        });
    }
}