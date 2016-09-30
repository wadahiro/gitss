import * as React from 'react';
import { Router, Route, Link, browserHistory } from 'react-router';

import Layout from './views/Layout';
import SearchView from './views/SearchView';
import NotFoundView from './views/NotFoundView';

export default class App extends React.Component<any, any> {
    render() {
        
        return ROUTES;
    }
}

export const ROUTES =
    <Router history={browserHistory}>
        <Route component={Layout}>
            <Route path="/" component={SearchView}/>
            <Route path="/issues" component={SearchView}/>
            <Route path="/issues/:_id" component={SearchView}/>
        </Route>
        <Route path="*" component={NotFoundView}/>
    </Router>;
