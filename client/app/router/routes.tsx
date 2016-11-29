import * as React from 'react';
import { Store } from 'redux';
import { Route, IndexRoute, Redirect } from 'react-router';

import { RootState } from '../reducers';
import HomeView from '../views/HomeView';
import SearchView from '../views/SearchView';
import NotFoundView from '../views/NotFoundView';

export default ({ store, first }: { store: Store<RootState>, first: { time: boolean } }) => {
    return (
        <Route>
            <Route path="/" component={HomeView} />
            <Route path="/search" component={SearchView} />
            <Route path="*" component={NotFoundView} />
        </Route>
    );
};