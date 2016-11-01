import * as React from 'react';

import { Tag } from './Tag';
import { Panel, PanelHeading, PanelBlock } from './Panel';
import { Menu, MenuLabel, MenuList, MenuLink } from './Menu';
import { Facet, FilterParams } from '../reducers';

interface FacetPanelProps {
    title: string;
    selected: string[];
    facet: Facet;
    onToggle: (term: string) => void;
    emptyKeyword?: string;
    emptyLabel?: string;
}

export function FacetPanel(props: FacetPanelProps) {
    if (!props.facet || !props.facet.terms) {
        return null;
    }

    return (
        <Panel>
            <PanelHeading>{props.title}</PanelHeading>
            <Menu>
                <MenuList>
                    {props.facet.terms.map(x => {
                        let type = x.term;
                        if (props.emptyKeyword && props.emptyLabel && type === props.emptyKeyword) {
                            type = props.emptyLabel;
                        }
                        const style = {
                            padding: 0
                        }
                        const isToggled = contains(props.selected, x.term);

                        return (
                            <PanelBlock key={type} style={style}>
                                <li onClick={props.onToggle.bind(null, x.term)}>
                                    <MenuLink count={x.count} isToggled={isToggled}>{type}</MenuLink>
                                </li>
                            </PanelBlock>
                        )
                    })}
                </MenuList>
            </Menu>
        </Panel>
    );
}

function contains(array = [], item) {
    return array !== null && typeof array.find(x => x === item) !== 'undefined';
}