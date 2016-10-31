import * as React from 'react';

import { Tag } from './Tag';
import { Panel, PanelHeading, PanelBlock } from './Panel';
import { Menu, MenuLabel, MenuList, MenuLink } from './Menu';
import { Facet } from '../reducers';

interface ExtFacetProps {
    facet: Facet;
}

export function ExtFacet(props: ExtFacetProps) {
    if (!props.facet || !props.facet.terms) {
        return null;
    }
    return (
        <Panel>
            <PanelHeading>File extensions</PanelHeading>
            <Menu>
                <MenuList>
                    {props.facet.terms.map(x => {
                        const type = x.term !== '' ? x.term : '(No extension)';
                        return (
                            <PanelBlock style={{ padding: 0 }}>
                                <li>
                                    <MenuLink href="#" count={x.count}>{type}</MenuLink>
                                </li>
                            </PanelBlock>
                        )
                    })}
                </MenuList>
            </Menu>
        </Panel>
    );

}