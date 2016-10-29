import * as React from 'react';

import { Tag } from './Tag';
import { Facets } from '../reducers';

const BMenu = require('re-bulma/lib/components/menu/menu').default;
const BMenuLabel = require('re-bulma/lib/components/menu/menu-label').default;
const BMenuList = require('re-bulma/lib/components/menu/menu-list').default;
const BMenuLink = require('re-bulma/lib/components/menu/menu-link').default;

interface FacetsProps {
    facets: Facets;
}

export function Facets(props: FacetsProps) {
    return (
        <BMenu {...props}>
            {Object.keys(props.facets).map(key => {
                const facet = props.facets[key];
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

interface MenuLinkProps {
    count: number;
    isActive?: boolean;
    children?: React.ReactElement<any>
}

export function MenuLink(props: MenuLinkProps) {
    return <BMenuLink {...props}>{props.children}<Tag size='isSmall' style={{ float: 'right' }}>{props.count}</Tag></BMenuLink>
}