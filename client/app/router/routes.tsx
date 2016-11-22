import * as React from 'react';
import { Route, IndexRoute, Redirect } from 'react-router';

import Layout from '../views/Layout';
import SearchView from '../views/SearchView';
import NotFoundView from '../views/NotFoundView';

/**
 * Returns configured routes for different
 * environments. `w` - wrapper that helps skip
 * data fetching with onEnter hook at first time.
 * @param {Object} - any data for static loaders and first-time-loading marker
 * @returns {Object} - configured routes
 */
export default ({store, first}) => {

    // Make a closure to skip first request
    function w(loader) {
        return (nextState, replaceState, callback) => {
            if (first.time) {
                first.time = false;
                return callback();
            }
            return loader ? loader({ store, nextState, replaceState, callback }) : callback();
        };
    }

    return (
        <Route component={Layout}>
            <Route path="/" component={SearchView} />
            <Route path="/s" component={SearchView} />
            <Route path="/s/:organization" component={SearchView} />
            <Route path="/s/:organization/:project" component={SearchView} />
            <Route path="/s/:organization/:project/:repository" component={SearchView} />
            <Route path="/s/:organization/:project/:repository/branches/:branch" component={SearchView} />
            <Route path="/s/:organization/:project/:repository/tags/:tag" component={SearchView} />
            <Route path="*" component={NotFoundView} />
        </Route>
    );

    // return <Route path="/" component={App}>
    //     <IndexRoute component={Homepage} onEnter={w(Homepage.onEnter)} />
    //     <Route path="/usage" component={Usage} onEnter={w(Usage.onEnter)} />
    //     {/* Server redirect in action */}
    //     <Redirect from="/docs" to="/usage" />
    //     <Route path="*" component={NotFound} onEnter={w(NotFound.onEnter)} />
    // </Route>;
};