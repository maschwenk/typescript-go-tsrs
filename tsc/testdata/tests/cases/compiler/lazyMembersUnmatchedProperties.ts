// @strict: true
// @target: esnext
// @noEmit: true

// Missing property errors and discriminant checks against instantiated class
// and interface references, including inherited members and the order in which
// missing properties are reported.

interface Base<T> {
    b1: T;
    b2: string;
    shared: number;
}
interface Derived<T> extends Base<T[]> {
    d1: T;
    shared: number;
    d2?: T;
}
const derived1: Derived<number> = {};
const derived2: Derived<number> = { d1: 1 };
const derived3: Derived<number> = { d1: 1, b1: [1], b2: "", shared: 1 };
class Impl<T> implements Derived<T> {
    d1!: T;
}

interface Many<T> extends Base<T> {
    m1: T;
    m2: T;
    m3: T;
    m4: T;
    m5?: T;
}
const many1: Many<string> = { m1: "" };
const many2: Many<string> = { b1: "", b2: "", shared: 0, m1: "", m2: "", m3: "", m4: "", extra: 1 };

class Box<T> {
    private secret!: T;
    contents!: T;
    method(): T {
        return this.contents;
    }
}
const box1: Box<string> = { contents: "" };
declare const notBox: { contents: string; method(): string };
const box2: Box<string> = notBox;

interface Circle<T> {
    kind: "circle";
    radius: T;
}
interface Square<T> {
    kind: "square";
    size: T;
}
type Shape<T> = Circle<T> | Square<T>;
const shape1: Shape<number> = { kind: "circle" };
const shape2: Shape<number> = { kind: "square", radius: 1 };
declare function area<T>(shape: Shape<T>): T;
const area1 = area({ kind: "circle", radius: 1 });
const area2 = area({ kind: "square", size: "x" });
declare function pick<T>(a: Circle<T> | T[]): T;
const picked = pick({ kind: "circle", radius: 1 });

interface Tree<T> extends Array<Tree<T>> {
    value: T;
}
const tree: Tree<string> = [];
