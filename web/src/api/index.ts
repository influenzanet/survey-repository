import { Store } from "../store";
import { Platform, SurveyStats } from "../types";

export interface ApiError {
    xhr: JQuery.jqXHR<any>
    text: string;
    error: Error
    url: string;
}

export const apiUrl = (path: string): string => {
    return  Store.api_url + '/' + path; 
}

export const apiload = <T>(url: string): Promise<T>=> {
    return new Promise((resolve, reject) => {

        const p = $.ajax({'url':url, 'dataType':'json'});

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
