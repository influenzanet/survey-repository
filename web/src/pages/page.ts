
import jQuery from "jquery";

export abstract class Page {

    name: string;

    constructor(name:string) {
        this.name = name;
    }

    abstract render():void;

    element(): JQuery {
        return jQuery('#page-'+ this.name);
    }

    show() {
        this.render();
        this.element().show();
    }

    forward(page: string) {
        $('html').trigger("app.page", page);
    }
}