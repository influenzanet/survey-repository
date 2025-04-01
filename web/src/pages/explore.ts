import {  apiLoadSurveys, SurveyFilters, apiLoadSurvey } from "../api";
import {  isSurveyGroupItem, LocalizedObjectBase, Survey, SurveyItem} from "survey-engine/data_types";
import { model_type_label, Paginated, SurveyMeta } from "../types";
import { h } from "../ui";
import { HtmlRenderer, HtmlRendererContext, BootstrapTheme   } from "ifn-survey-tools"
import { Page } from "./page";

const from_timestamp = (time:number)=> {
    const d = new Date(time * 1000);
    return d.toISOString();
}


export class Explore extends Page {

    constructor() {
        super('explore');
    }

    render(): void {
        const filters: SurveyFilters = {
            limit: 100,
            offset: 0,
        };
        apiLoadSurveys('influenzanet', filters).then(results => {
            this.buildTable(results);
        });
    }

    buildTable(results: Paginated<SurveyMeta>) {
        const $ui = $('#explore-ui');

        $ui.empty();

        const $tb = h('table', {"class": "table"});

        const headers : string[] = [
            'Id',
            'Platform',
            'Type',
            'Version',
            'Name',
            'Published'
        ];

        const $h = h('thead');

        headers.forEach(v => {
            $h.append(h('th', {}, v));
        });
        $tb.append($h);

        const $tbody = h('tbody');
        $tb.append($tbody);
        results.data.forEach(row => {
            const $r = h('tr');

            const td = (value: string)=> {
                return h('td', {}, value);
            }

            $r.append( td('' + row.id) );
            $r.append( td(row.platform));
            $r.append( td(row.model_type) );
            $r.append( td(row.version));
            $r.append( td(row.descriptor.name));
            $r.append( td(from_timestamp(row.imported_at)));

            const $btnShow = h('button', {"class":"btn btn-sm btn-info"}, "Show");

            $btnShow.on('click', ()=> {
                this.showSurvey(row);
            });

            $r.append( h('td', {}, $btnShow));

            $tbody.append($r);
        });

        $ui.append($tb);
    }

    showSurvey(meta: SurveyMeta) {
        apiLoadSurvey(meta.id).then(data=> {
            let content = '';
            if(meta.model_type == 'D') {
                const survey = data as Survey;
                content = render_survey_definition(survey);
            } else {
                content = '<pre>' + JSON.stringify(data, undefined, 2) + '</pre>';
            }
            const $explore =  $('#explore-ui');
            const $ui = $('#survey-show');
            const $f = $('<iframe style="width:100%">');
            $f.attr('srcdoc', content);
            $ui.empty();
            
            const $btn = h('button', {"class": "btn btn-sm btn-info"}, "Return to table");
            $btn.on('click', ()=> {
                $ui.hide();
                $explore.show();
            });

            $ui.append($btn);

            $ui.append( h('h3', {}, 'Survey ' + meta.id + ' Type=' + model_type_label(meta.model_type) +' Version=' + meta.version + ' Platform=' + meta.platform + ' Name=' + meta.descriptor.name));

            $ui.show();
            $explore.hide();
            $ui.append($f);
        })
    }
}


const findSurveyLanguages = (survey: Survey):Set<string> =>{
    const lang = new Set<string>();
    const extractCode = (loc: LocalizedObjectBase[]|undefined) =>{
        if(!loc) {
            return;
        }
        loc.forEach(o=> lang.add(o.code));
    }
    
    var itemCount = 0;
    
    const visitItem = (item:SurveyItem) =>{
        if(itemCount > 20) {
            return;
        }
        if(isSurveyGroupItem(item)) {
            item.items.forEach(visitItem);
        } else {
            if(item.components) {
                let has = false;
                item.components.items.forEach(comp=> {
                    extractCode(comp.content);
                    extractCode(comp.description);
                    has = true;
                });
                if(has) {
                    itemCount += 1;
                }
            }
        }
    }

    if(survey.props) {
        extractCode(survey.props.description);
        extractCode(survey.props.name);
        extractCode(survey.props.typicalDuration);
    } else {
        survey.surveyDefinition.items.forEach(visitItem);
    }

    return lang;
}


const render_survey_definition = (survey: Survey):string => {
    
    const languages = findSurveyLanguages(survey);

    const lang = Array.from(languages.values());

    // Create a survey context, to tell the renderer how to render (languages to show and css theme)
    const context = new HtmlRendererContext({languages: lang}, new BootstrapTheme());
    const renderer = new HtmlRenderer();
    
    // Considering `survey` contains your survey definition
    return renderer.render(survey, context);
}