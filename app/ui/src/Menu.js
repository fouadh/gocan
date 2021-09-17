import {Menubar} from "primereact/menubar";

export function Menu() {
    const items = [
        {
            label: "Gocan UI",
            icon: "pi pi-chart-bar",
            url: "/"
        }
    ];

    return <Menubar model={items}/>;
}