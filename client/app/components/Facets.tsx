import * as React from 'react';

import { Tag } from './Tag';
import { Facets as F } from '../reducers';

const BMenu = require('re-bulma/lib/components/menu/menu').default;
const BMenuLabel = require('re-bulma/lib/components/menu/menu-label').default;
const BMenuList = require('re-bulma/lib/components/menu/menu-list').default;
const BMenuLink = require('re-bulma/lib/components/menu/menu-link').default;

interface FacetsProps {
    facets: F;
}

export class Facets extends React.PureComponent<FacetsProps, void> {
    render() {
        return (
            <BMenu {...this.props}>
                {Object.keys(this.props.facets).map(key => {
                    const facet = this.props.facets[key];
                    return (
                        <div key={key}>
                            <BMenuLabel>{facet.field}</BMenuLabel>
                            <BMenuList>
                                {facet.terms.map(term => {
                                    return (
                                        <li key={term.term}><MenuLink count={term.count}>{term.term}</MenuLink></li>
                                    );
                                })}
                            </BMenuList>
                        </div>
                    );
                })}
            </BMenu>
        );
    }
}

interface MenuLinkProps {
    count: number;
    isActive?: boolean;
    children?: React.ReactElement<any>
}

export class MenuLink extends React.PureComponent<MenuLinkProps, void> {
    render() {
        return <BMenuLink {...this.props}>{this.props.children}<Tag size='isSmall' style={{ float: 'right' }}>{this.props.count}</Tag></BMenuLink>;
    }
}