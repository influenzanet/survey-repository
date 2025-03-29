
type ModelType = 'D' | 'P';

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