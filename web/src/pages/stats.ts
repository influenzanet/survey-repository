import { apiLoadStats } from "../api";
import { SurveyStats } from "../types";
import { h, icon } from "../ui";
import { Page } from "./page";

interface SurveyStat {
    definition: number;
    preview: number;
}

const empty_stat = ()=> { 
    return {"definition": 0, "preview": 0};
}

export class Stats extends Page {

    constructor() {
        super('stats');
    }

    render(): void {
        apiLoadStats('influenzanet').then(stats => {
            this.buildTable(stats);
        });
    }

    buildTable(stats: SurveyStats[]) {
        const surveys = new Set<string>();
            const platforms: Map<string, Map<string, SurveyStat>> = new Map();
            stats.forEach((r)=> {
                surveys.add(r.survey_key);
                let p = platforms.get(r.platform);
                if(!p) {
                    p = new Map();
                    platforms.set(r.platform, p);
                }
                let s = p.get(r.survey_key);
                if(!s) {
                    s = empty_stat();
                    p.set(r.survey_key, s);
                }
                switch(r.model_type) {
                    case 'D':
                        s.definition = r.count;
                        break;
                    case 'P':
                        s.preview = r.count;
                    break;                
                }
            });
            const $tb = $('<table class="table text-center">');
          
            const $head = $('<thead>');
            $tb.append($head);
            
            const $h = $('<tr>');
            $h.append( h('th', {}, "Platform"));
            $head.append($h);
            
            surveys.forEach(survey=> {
                $h.append( h('th', undefined, survey) );
            });

            const $body = $('<tbody>');
            $tb.append($body);
            
            platforms.forEach((ss, platform_code)=> {
                const $row = $('<tr>');
                $body.append($row);

                $row.append(platform_code);
                surveys.forEach(survey=> {
                    const s = ss.get(survey) ?? empty_stat();
                    const $c = $('<td>');
                    $c.append(h('span', {"class":"badge bg-danger me-2"} , icon('sitemap') + ' ' + s.definition));
                    $c.append(h('span', {"class":"badge bg-warning"} , icon('rectangle-list') + ' ' + s.preview));
                    $row.append($c);
                });
            });

            const $ui = $('#stats-ui');
            $ui.empty();
            $ui.append($tb);
    }
}