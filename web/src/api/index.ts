import { Store } from "../store";
import { Paginated, Platform, SurveyMeta, SurveyStats } from "../types";

export interface ApiError {
    xhr: JQuery.jqXHR<any>
    text: string;
    error: Error
    url: string;
}

export const apiUrl = (path: string): string => {
    return  Store.api_url + '/' + path; 
}

export const apiload = <T>(url: string, params?: object): Promise<T>=> {
    return new Promise((resolve, reject) => {
        let u: string;
        if(params) {
            u = url + '?' + jQuery.param(params);
        } else {
            u = url;
        }

        const p = $.ajax({'url':u, 'dataType':'json'});

        p.then((json) => {
            resolve(json);
        });

        p.catch( (jxhr,text, error) =>{
            console.error(jxhr, text, error);
            reject({'xhr': jxhr, 'text': text, 'error': error, 'url': url });
        })
    });
}

export const apiLoadStats = (namespace:string): Promise<SurveyStats[]> => {
    return apiload<SurveyStats[]>(apiUrl('namespace/'+namespace+'/surveys/stats'))
}

export const apiLoadPlatforms = (): Promise<Platform[]> => {
    return apiload<Platform[]>(apiUrl('refs/platforms'))
}

export interface SurveyFilters {
    limit: number
    offset: number
}


export const apiLoadSurveys = (namespace:string, filter?: SurveyFilters ): Promise<Paginated<SurveyMeta>> => {
    const params: Record<string, any> = {};
    if(filter) {
        if(filter.limit) {
            params['limit'] = filter.limit;
            if(filter.offset) {
                params['offset'] = filter.offset;
            }
        }
    }
    return apiload<Paginated<SurveyMeta>>(apiUrl('namespace/'+namespace+'/surveys'), params)
}

export const apiLoadSurvey = (id: number): Promise<object> =>{
    return apiload(apiUrl('survey/'+id +'/data'));
}

