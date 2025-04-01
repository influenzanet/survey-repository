
type ModelType = 'D' | 'P';


export const model_type_label= (model_type: string) => {
    switch(model_type) {
        case 'P':
            return 'Preview';
        case 'D':
            return 'Definition';
        default:
            return 'Unknown';
    }
}



export interface Platform {
    id: string;
    label: string;
    country: string;
}

export interface SurveyStats {
    platform: string
    model_type:ModelType
    survey_key:string
    count: number;
}

export interface SurveyDescriptor {
    name: string;
    version: string;
    external_id: string;
    published: number;
    model_version: string;
    sha256: string;
}

export interface SurveyMeta {
    id: number;
    namespace: string;
    imported_at: number;
    imported_by: string;
    platform: string;
    version: string;
    model_type: string;
    labels: Record<string,string>;
    descriptor: SurveyDescriptor
}

export interface Paginated<T> {
    offset: number;
    limit: number;
    total: number;
    data: T[]
}
