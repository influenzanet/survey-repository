import { Explore } from "./explore";
import { Home } from "./home";
import { Page } from "./page";
import { Platform } from "./platforms";
import { Stats } from "./stats";

export const createPage = (page:string) => {
    var p: Page|undefined = undefined;
    switch(page) {
        case 'home':
            p = new Home();
            break;
        case 'stats':
            p = new Stats();
            break;
        case 'platforms':
            p = new Platform();
            break;
        case 'explore':
            p = new Explore();
            break;
    }
    if(p) {
        p.show();
    } else {
        console.error("Unknown page "+ page);
    }
}