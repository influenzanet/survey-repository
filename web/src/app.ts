import { createPage } from "./pages";
import { Store, StoreType } from "./store";
import jQuery from "jquery";
import {defer} from "lodash";
import { showError } from "./ui";
const $ = jQuery;

export interface AppInterface {
    store: StoreType
    init():void
    start():void
    load():Promise<boolean>
    show(page:string):void
    showError(error: any):void
}

export const App: AppInterface = {
    store: Store,
    
    init: function() {
        $('.page').hide(); // Hide all page content at start
        let u: string = '';
        if(import.meta.env.DEV) {
            console.log('Dev mode');
            u = import.meta.env.VITE_API_URL;
            if(typeof(u)=="undefined") {
                console.warn("evn VITE_API_URL might not be set !");
                u = '';
            }
        } else {
            u = $('html').data('url') ?? '';
        }
        Store.api_url = u;
        console.log('api=', Store.api_url);
        $('.nav-page').on('click', function() {
            const target = $(this).data('target');
            defer(()=>{
                App.show(target);
            });
        });
    },

    start: function() {
        console.log('start');
        App.init();
        App.load().then(()=> {
            const p = document.location.pathname.slice(1);
            const page = p != '' ? p : 'home';
            console.log('go to page', page);
            App.show(page);
        }).catch((reason)=>{
            App.showError(reason);
        });
    },

    load: function() {
        return new Promise<boolean>((resolve, _) => {
            resolve(true);
        });
    },

    show: function(page:string) {
        $('.page').hide();
        
        const p = page.split('/');

        const name = p[0];
        //const args = p.length > 1 ? p.slice(1) : undefined;
        createPage(name);
    },
    showError: function(error: any) {
        showError(error);
    }

};

window.App = App;

