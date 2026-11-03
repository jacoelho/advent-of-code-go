#!/usr/bin/env python3
"""Independent stdlib AoC2022 oracles; no Go outputs are used.

Usage: python3 migration/oracles/2022.py [DAY ...] [--input-dir PATH] [--output PATH]
Default inputs are repository inputs/2022/NN.txt. Source dayNN.txt names also work.
Each day has a 90-second elapsed-time limit; results are checkpointed per day.
"""
import argparse, collections, functools, hashlib, heapq, itertools, json, math
import pathlib, re, signal, time
from fractions import Fraction

ROOT=pathlib.Path(__file__).resolve().parent.parent.parent/'inputs'/'2022'
REV='c420b42f08c6e8687d19f7e446c04b3feab88b04'
SOURCE=[(69177,207456),(11475,16862),(8233,2821),(569,936),('FZCMJCRHZ','JSDHQMZGF'),
        (1578,2178),(1427048,2940614),(1801,209880),(6269,2557),(14240,'PLULKBZH'),
        (99840,20683044837),(447,446),(5393,26712),(1513,22646),(5525847,13340867187704),
        (1751,2207),(3059,1500874635587),(4340,2468),(1487,13440),(23321,1428396909280),
        (223971851179174,3379022190351),(88226,57305),(3762,997),(322,974),('2011-=2=-1020-1===-1',None)]

def d01(t):
    sums=sorted((sum(map(int,b.split())) for b in t.strip().split('\n\n')),reverse=True)
    return sums[0],sum(sums[:3]),{'method':'explicit per-elf totals and sorting'}
def d02(t):
    a=b=0
    for row in t.splitlines():
        x,y=row.split();x=ord(x)-65;y=ord(y)-88
        a+=y+1+(3 if x==y else 6 if (y-x)%3==1 else 0)
        b+=(x+y-1)%3+1+3*y
    return a,b,{'method':'modular cyclic move relation and explicit requested-outcome score'}
def d03(t):
    rows=t.split();assert len(rows)%3==0
    def priority(items):
        assert len(items)==1;c=next(iter(items));return ord(c)-ord('a')+1 if c.islower() else ord(c)-ord('A')+27
    a=0
    for s in rows:
        assert len(s)%2==0;n=len(s)//2;a+=priority(set(s[:n])&set(s[n:]))
    b=sum(priority(set(rows[i])&set(rows[i+1])&set(rows[i+2])) for i in range(0,len(rows),3))
    return a,b,{'method':'literal compartment/group set intersections, verify unique common type'}
def d04(t):
    contained=overlap=0
    for row in t.splitlines():
        a,b,c,d=map(int,re.findall(r'\d+',row));contained+=(a<=c<=d<=b or c<=a<=b<=d);overlap+=max(a,c)<=min(b,d)
    return contained,overlap,{'method':'inclusive interval containment/intersection'}
def d05(t):
    drawing,moves=t.split('\n\n',1);rows=drawing.splitlines();labels=list(re.finditer(r'\d+',rows.pop()));stacks=[]
    for label in labels:
        col=label.end()-1;stacks.append([row[col] for row in reversed(rows) if col<len(row) and row[col]!=' '])
    def solve(reverse):
        s=[r[:] for r in stacks]
        for row in moves.splitlines():
            if not row:continue
            n,a,b=map(int,re.fullmatch(r'move (\d+) from (\d+) to (\d+)',row).groups());a-=1;b-=1;assert n<=len(s[a])
            if n==0 or a==b:continue
            block=s[a][-n:];del s[a][-n:];s[b].extend(reversed(block) if reverse else block)
        return ''.join(s[-1] for s in s)
    return solve(True),solve(False),{'method':'parse ORIGINAL drawing columns; independent copied stack simulation','initial_stacks_bottom_to_top':[''.join(s) for s in stacks]}
def d06(t):
    s=t.strip('\r\n');assert '\n' not in s
    def find(n):return next(i for i in range(n,len(s)+1) if len(set(s[i-n:i]))==n)
    return find(4),find(14),{'method':'direct independent distinct-window scan, terminator excluded'}
def d07(t):
    cwd=();dirs={()};files={}
    for row in t.splitlines():
        p=row.split()
        if p[:2]==['$','cd']:
            name=p[2];cwd=() if name=='/' else cwd[:-1] if name=='..' else cwd+(name,);dirs.add(cwd)
        elif p[:2]==['$','ls']:continue
        elif p[0]=='dir':dirs.add(cwd+(p[1],))
        else:
            key=cwd+(p[1],);v=int(p[0]);assert key not in files or files[key]==v;files[key]=v
    sizes={d:sum(v for p,v in files.items() if p[:len(d)]==d) for d in dirs};needed=max(0,sizes[()]-40000000)
    return sum(v for v in sizes.values() if v<=100000),min(v for v in sizes.values() if v>=needed) if needed else 0,{'method':'canonical file identities and prefix-derived directory totals','used_space':sizes[()]}
def d08(t):
    a=[[int(c) for c in row] for row in t.splitlines()];h=len(a);w=len(a[0]);visible=best=0
    for y in range(h):
        for x in range(w):
            rays=([a[y][u] for u in range(x-1,-1,-1)],[a[y][u] for u in range(x+1,w)],[a[v][x] for v in range(y-1,-1,-1)],[a[v][x] for v in range(y+1,h)])
            visible+=any(all(v<a[y][x] for v in ray) for ray in rays);score=1
            for ray in rays:
                n=0
                for v in ray:
                    n+=1
                    if v>=a[y][x]:break
                score*=n
            best=max(best,score)
    return visible,best,{'method':'literal four directional sight rays including blocking tree'}
def d09(t):
    moves=[(s.split()[0],int(s.split()[1])) for s in t.splitlines()];delta={'R':(1,0),'L':(-1,0),'U':(0,1),'D':(0,-1)}
    def solve(n):
        rope=[(0,0)]*n;seen={(0,0)}
        for op,k in moves:
            dx,dy=delta[op]
            for _ in range(k):
                rope[0]=(rope[0][0]+dx,rope[0][1]+dy)
                for i in range(1,n):
                    a,b=rope[i-1];x,y=rope[i];u,v=a-x,b-y
                    if max(abs(u),abs(v))>1:rope[i]=(x+(u>0)-(u<0),y+(v>0)-(v<0))
                seen.add(rope[-1])
        return len(seen)
    return solve(2),solve(10),{'method':'literal signed knot positions, initial origin included'}
def crt(t):
    x=1;cycles=[]
    for row in t.splitlines():
        p=row.split();cycles.append(x)
        if p[0]=='addx':cycles.append(x);x+=int(p[1])
        else:assert p==['noop']
    assert len(cycles)>=240
    signal=sum(c*cycles[c-1] for c in (20,60,100,140,180,220))
    rows=[''.join('#' if abs(col-cycles[r*40+col])<=1 else '.' for col in range(40)) for r in range(6)]
    return signal,rows,len(cycles)
def d10(t):
    signal,rows,cycles=crt(t)
    font={'###./#..#/#..#/###./#.../#...':'P','#.../#.../#.../#.../#.../####':'L',
          '#..#/#..#/#..#/#..#/#..#/.##.':'U','#..#/#.#./##../#.#./#.#./#..#':'K',
          '###./#..#/###./#..#/#..#/###.':'B','####/...#/..#./.#../#.../####':'Z',
          '#..#/#..#/####/#..#/#..#/#..#':'H'}
    text=''.join(font['/'.join(r[i:i+4] for r in rows)] for i in range(0,40,5))
    return signal,text,{'method':'literal per-cycle X timeline, 240 independent pixels and standard glyph transcription','rendered_rows':rows,'cycles':cycles}
def d11(t):
    monkeys={}
    for b in t.strip().split('\n\n'):
        r=[s.strip() for s in b.splitlines()];i=int(re.search(r'\d+',r[0])[0]);items=list(map(int,re.findall(r'\d+',r[1])));op=re.fullmatch(r'Operation: new = old ([+*]) (old|\d+)',r[2]);divisor=int(re.search(r'\d+',r[3])[0]);yes=int(re.search(r'\d+',r[4])[0]);no=int(re.search(r'\d+',r[5])[0]);monkeys[i]=(items,op[1],op[2],divisor,yes,no)
    ids=sorted(monkeys);mod=math.lcm(*(m[3] for m in monkeys.values()))
    def solve(rounds,relief):
        queues={i:collections.deque(monkeys[i][0]) for i in ids};counts=collections.Counter()
        for _ in range(rounds):
            for i in ids:
                _,op,arg,d,yes,no=monkeys[i];assert yes!=i and no!=i
                while queues[i]:
                    v=queues[i].popleft();b=v if arg=='old' else int(arg);v=v+b if op=='+' else v*b;v=v//3 if relief else v%mod;queues[yes if v%d==0 else no].append(v);counts[i]+=1
        return math.prod(sorted(counts.values())[-2:]),dict(counts)
    a,p=solve(20,True);b,q=solve(10000,False)
    return a,b,{'method':'exact-integer queues in ID order; P1 divide3, P2 LCM residues','inspection_counts':[p,q],'lcm':mod}
def d12(t):
    rows=t.splitlines();h=len(rows);w=len(rows[0]);height={};start=end=None
    for y,row in enumerate(rows):
        for x,c in enumerate(row):
            if c=='S':start=(x,y);c='a'
            if c=='E':end=(x,y);c='z'
            height[x,y]=ord(c)-97
    q=collections.deque([end]);dist={end:0}
    while q:
        x,y=q.popleft()
        for p in ((x-1,y),(x+1,y),(x,y-1),(x,y+1)):
            if p in height and p not in dist and height[p]>=height[x,y]-1:dist[p]=dist[x,y]+1;q.append(p)
    return dist[start],min(v for p,v in dist.items() if height[p]==0),{'method':'single exact reverse BFS from E'}
def packetcmp(a,b):
    if isinstance(a,int) and isinstance(b,int):return (a>b)-(a<b)
    if isinstance(a,int):a=[a]
    if isinstance(b,int):b=[b]
    for x,y in zip(a,b):
        c=packetcmp(x,y)
        if c:return c
    return (len(a)>len(b))-(len(a)<len(b))
def d13(t):
    pairs=[list(map(json.loads,b.splitlines())) for b in t.strip().split('\n\n')];a=sum(i+1 for i,(x,y) in enumerate(pairs) if packetcmp(x,y)<0);packets=[(p,None) for pair in pairs for p in pair]+[([[2]],2),([[6]],6)];packets.sort(key=functools.cmp_to_key(lambda x,y:packetcmp(x[0],y[0])));where=[i+1 for i,p in enumerate(packets) if p[1] is not None]
    return a,math.prod(where),{'method':'exact JSON trees, recursive promoting comparison; inserted-divider identity','divider_positions':where}
def d14(t):
    rocks=set()
    for row in t.splitlines():
        pts=[tuple(map(int,p.split(','))) for p in row.split(' -> ')]
        for (x,y),(u,v) in zip(pts,pts[1:]):
            assert x==u or y==v
            for a in range(min(x,u),max(x,u)+1):
                for b in range(min(y,v),max(y,v)+1):rocks.add((a,b))
    bottom=max(y for x,y in rocks)
    def solve(floor):
        occupied=set(rocks);n=0
        while (500,0) not in occupied:
            x,y=500,0
            while True:
                if not floor and y>bottom:return n
                options=[(x,y+1),(x-1,y+1),(x+1,y+1)];p=next((p for p in options if p not in occupied and (not floor or p[1]!=bottom+2)),None)
                if p is None:occupied.add((x,y));n+=1;break
                x,y=p
        return n
    return solve(False),solve(True),{'method':'literal individual-grain simulation; floor represented by y only'}

def merge(intervals):
    result=[]
    for a,b in sorted(intervals):
        if result and a<=result[-1][1]+1:result[-1]=(result[-1][0],max(b,result[-1][1]))
        else:result.append((a,b))
    return result
def triangular_domain_count(k,lo,hi):
    """Count 0<=w<=k, lo<=v<=hi, |v|<=w, w==v mod2."""
    if k<0:return 0
    def segment(a,b):
        a=max(0,a);b=min(k,b)
        if a>b:return 0
        def prefix(n):return n*n//4 if n>=0 else 0
        return prefix(k-a)-prefix(k-b-1)+b-a+1
    return segment(max(0,lo),hi)+segment(max(1,-hi),-lo)
def sensor_results(t,row,limit):
    sensors=[];beacons=set()
    for line in t.splitlines():
        x,y,bx,by=map(int,re.findall(r'-?\d+',line));sensors.append((x,y,abs(x-bx)+abs(y-by)));beacons.add((bx,by))
    intervals=merge((x-(r-abs(y-row)),x+(r-abs(y-row))) for x,y,r in sensors if abs(y-row)<=r)
    covered=sum(b-a+1 for a,b in intervals)-sum(y==row and any(a<=x<=b for a,b in intervals) for x,y in beacons)
    events=collections.defaultdict(list);events[0];events[2*limit+1];events[limit+1]
    for i,(x,y,r) in enumerate(sensors):
        lo=max(0,x+y-r);hi=min(2*limit,x+y+r)
        if lo<=hi:events[lo].append((i,True));events[hi+1].append((i,False))
    active=set();breaks=sorted(events);total=0;point=None;slabs=0
    for k,u0 in enumerate(breaks[:-1]):
        for i,add in events[u0]:
            if add:active.add(i)
            else:active.remove(i)
        u1=breaks[k+1]-1
        if u0>u1:continue
        slabs+=1;coveredv=merge((max(-limit,sensors[i][0]-sensors[i][1]-sensors[i][2]),min(limit,sensors[i][0]-sensors[i][1]+sensors[i][2])) for i in active if sensors[i][0]-sensors[i][1]-sensors[i][2]<=limit and sensors[i][0]-sensors[i][1]+sensors[i][2]>=-limit)
        gaps=[];v=-limit
        for a,b in coveredv:
            if v<a:gaps.append((v,a-1))
            v=max(v,b+1)
        if v<=limit:gaps.append((v,limit))
        left,right=(u0,u1) if u1<=limit else (2*limit-u1,2*limit-u0)
        for lo,hi in gaps:
            n=triangular_domain_count(right,lo,hi)-triangular_domain_count(left-1,lo,hi)
            if n:
                a,b=left,right
                while a<b:
                    mid=(a+b)//2
                    if triangular_domain_count(mid,lo,hi)>triangular_domain_count(left-1,lo,hi):b=mid
                    else:a=mid+1
                w=a;v=max(lo,-w);v+=(w-v)%2;u=w if u1<=limit else 2*limit-w;point=((u+v)//2,(u-v)//2)
            total+=n
    assert total==1,('uncovered count',total)
    return covered,point[0]*4000000+point[1],{'method':'whole-square exact transformed-lattice slab sweep; closed-form uncovered-point count','row':row,'limit':limit,'uncovered_count':total,'uncovered_point':point,'sweep_slabs':slabs}
def d15(t):return sensor_results(t,2000000,4000000)

def d16(t):
    flows={};graph={}
    for row in t.splitlines():
        names=re.findall(r'[A-Z]{2}',row);a=names[0];flows[a]=int(re.search(r'rate=(\d+)',row)[1]);graph[a]=names[1:]
    useful=sorted(a for a in flows if flows[a]>0);sources=['AA']+[a for a in useful if a!='AA'];distance={}
    for a in sources:
        q=collections.deque([a]);ds={a:0}
        while q:
            u=q.popleft()
            for v in graph[u]:
                if v not in ds:ds[v]=ds[u]+1;q.append(v)
        distance[a]=ds
    def enumerate_scores(minutes):
        bestmask={0:0};seen={}
        def visit(pos,left,opened,score):
            key=(pos,left,opened)
            if seen.get(key,-1)>=score:return
            seen[key]=score;bestmask[opened]=max(score,bestmask.get(opened,0))
            for i,v in enumerate(useful):
                if opened>>i&1 or v not in distance[pos]:continue
                rem=left-distance[pos][v]-1
                if rem>0:visit(v,rem,opened|(1<<i),score+rem*flows[v])
        visit('AA',minutes,0,0);return bestmask,len(seen)
    first,n=enumerate_scores(30);second,m=enumerate_scores(26);size=1<<len(useful);subset=[second.get(i,0) for i in range(size)]
    for i in range(len(useful)):
        for mask in range(size):
            if mask>>i&1:subset[mask]=max(subset[mask],subset[mask^(1<<i)])
    best=max(score+subset[(size-1)^mask] for mask,score in second.items())
    return max(first.values()),best,{'method':'exact location/time/opened-mask dominance; exhaustive reachable mask scores and disjoint-subset maxima','useful_valves':useful,'dominance_states':[n,m]}

SHAPES=(((0,0),(1,0),(2,0),(3,0)),((1,0),(0,1),(1,1),(2,1),(1,2)),((0,0),(1,0),(2,0),(2,1),(2,2)),((0,0),(0,1),(0,2),(0,3)),((0,0),(1,0),(0,1),(1,1)))
def rocks(jets,count,skip=True):
    occupied={(x,0) for x in range(7)};height=0;jet=0;number=0;added=0;seen={};proof=None
    while number<count:
        shape=SHAPES[number%5];x=2;y=height+4
        while True:
            dx=1 if jets[jet]=='>' else -1;jet=(jet+1)%len(jets)
            if all(0<=x+dx+a<7 and (x+dx+a,y+b) not in occupied for a,b in shape):x+=dx
            if all((x+a,y-1+b) not in occupied for a,b in shape):y-=1
            else:break
        occupied.update((x+a,y+b) for a,b in shape);height=max(height,max(y+b for a,b in shape));number+=1
        # Only rocks adjacent to top-connected empty air can affect future collisions.
        air={(x,height+1) for x in range(7)};q=list(air);surface=set()
        for a,b in q:
            for p in ((a-1,b),(a+1,b),(a,b-1),(a,b+1)):
                u,v=p
                if not 0<=u<7 or not 0<=v<=height+1:continue
                if p in occupied:surface.add(p)
                elif v>=1 and p not in air:air.add(p);q.append(p)
        occupied=surface
        if skip and proof is None:
            key=(number%5,jet,tuple(sorted((a,b-height) for a,b in occupied)))
            if key in seen:
                oldn,oldh=seen[key];period=number-oldn;gain=height-oldh;repeats=(count-number)//period
                proof={'first_rocks':oldn,'repeat_rocks':number,'first_height':oldh,'repeat_height':height,'period':period,'height_gain':gain,'next_shape':number%5,'next_jet':jet,'retained_rocks':len(occupied),'normalized_surface':key[2],'state_sha256':hashlib.sha256(repr(key).encode()).hexdigest(),'skipped_cycles':repeats}
                number+=repeats*period;added+=repeats*gain
            else:seen[key]=(number,height)
    return height+added,proof
def d17(t):
    jets=t.strip();a,_=rocks(jets,2022,False);b,proof=rocks(jets,1000000000000)
    assert proof is not None
    return a,b,{'method':'literal2022; exact top-reachable-air boundary/shape/jet cycle, normalized translation','cycle':proof}

NEIGHBORS3=((1,0,0),(-1,0,0),(0,1,0),(0,-1,0),(0,0,1),(0,0,-1))
def d18(t):
    cubes={tuple(map(int,row.split(','))) for row in t.splitlines()};lo=tuple(min(c[i] for c in cubes)-1 for i in range(3));hi=tuple(max(c[i] for c in cubes)+1 for i in range(3));air={lo};q=collections.deque([lo]);outside=0
    while q:
        c=q.popleft()
        for d in NEIGHBORS3:
            p=tuple(c[i]+d[i] for i in range(3))
            if p in cubes:outside+=1
            elif all(lo[i]<=p[i]<=hi[i] for i in range(3)) and p not in air:air.add(p);q.append(p)
    area=sum(tuple(c[i]+d[i] for i in range(3)) not in cubes for c in cubes for d in NEIGHBORS3)
    return area,outside,{'method':'literal exposed faces plus independently bounded exterior-air flood','air_cells':len(air)}

def blueprint_max(costs,minutes):
    caps=tuple(max(c[i] for c in costs) for i in range(3))
    @functools.cache
    def best(left,ore,clay,obs,ro,rc,rb,rg):
        resources=(ore,clay,obs);robots=(ro,rc,rb,rg);answer=rg*left
        for kind in (3,2,1,0):
            if kind<3 and robots[kind]>=caps[kind]:continue
            wait=0
            for i in range(3):
                need=costs[kind][i]-resources[i]
                if need>0:
                    if robots[i]==0:break
                    wait=max(wait,(need+robots[i]-1)//robots[i])
            else:
                elapsed=wait+1;rem=left-elapsed
                if rem<=0:continue
                nr=list(robots);nr[kind]+=1
                if rg*elapsed+nr[3]*rem+rem*(rem-1)//2<=answer:continue
                stock=[resources[i]+robots[i]*elapsed-costs[kind][i] for i in range(3)]
                stock=[min(stock[i],max(0,caps[i]*rem-nr[i]*(rem-1))) for i in range(3)]
                value=rg*elapsed+best(rem,*stock,*nr);answer=max(answer,value)
        return answer
    value=best(minutes,0,0,0,1,0,0,0);return value,best.cache_info().currsize
def d19(t):
    blueprints=[]
    for row in t.splitlines():
        i,a,b,c,d,e,f=map(int,re.findall(r'\d+',row));blueprints.append((i,((a,0,0),(b,0,0),(c,d,0),(e,0,f))))
    first=[];second=[];counts=[]
    for i,cost in blueprints:
        v,n=blueprint_max(cost,24);first.append((i,v));counts.append(n)
    for i,cost in blueprints[:3]:
        v,n=blueprint_max(cost,32);second.append((i,v));counts.append(n)
    return sum(i*v for i,v in first),math.prod(v for i,v in second),{'method':'exact wait-to-afford memo search, useful-resource caps, admissible triangular future-geode upper bound','geodes24':first,'geodes32':second,'memo_states':counts}

def d20(t):
    values=list(map(int,t.split()));assert values.count(0)==1 and len(values)>1;n=len(values)
    def solve(key,rounds):
        scaled=[v*key for v in values];order=list(range(n))
        for _ in range(rounds):
            for i,v in enumerate(scaled):
                pos=order.index(i);order.pop(pos);order.insert((pos+v)%(n-1),i)
        zero=order.index(values.index(0));return sum(scaled[order[(zero+i)%n]] for i in (1000,2000,3000))
    return solve(1,1),solve(811589153,10),{'method':'direct stable-identity circular list mixing using modulo n-1'}

def d21(t):
    nodes={}
    for row in t.splitlines():
        name,expr=row.split(': ');p=expr.split();nodes[name]=int(p[0]) if len(p)==1 else tuple(p)
    @functools.cache
    def evalnode(name,human=None):
        if name=='humn' and human is not None:return human
        v=nodes[name]
        if isinstance(v,int):return Fraction(v)
        a,op,b=v;x,y=evalnode(a,human),evalnode(b,human)
        return {'+':lambda:x+y,'-':lambda:x-y,'*':lambda:x*y,'/':lambda:x/y}[op]()
    @functools.cache
    def affine(name):
        if name=='humn':return Fraction(1),Fraction(0)
        v=nodes[name]
        if isinstance(v,int):return Fraction(0),Fraction(v)
        a,op,b=v;aa,ab=affine(a);ba,bb=affine(b)
        if op=='+':return aa+ba,ab+bb
        if op=='-':return aa-ba,ab-bb
        if op=='*':
            assert not aa or not ba,'nonlinear multiplication';return aa*bb+ba*ab,ab*bb
        assert not ba,'variable denominator';return aa/bb,ab/bb
    first=evalnode('root');assert first.denominator==1
    left,_,right=nodes['root'];a,b=affine(left);c,d=affine(right);assert a!=c;human=(d-b)/(a-c);assert human.denominator==1
    lv,rv=evalnode(left,human),evalnode(right,human);assert lv==rv
    assert all(evalnode(name,h).denominator==1 for h in (None,human) for name in nodes),'nonintegral literal intermediate'
    return int(first),int(human),{'method':'independent exact-rational affine symbolic equation, literal substitution','root_coefficients':[[str(a),str(b)],[str(c),str(d)]],'verified_root_sides':[str(lv),str(rv)],'literal_integer_intermediates_verified':True}

def neg(v):return tuple(-x for x in v)
def fold_orient(orientation,direction):
    u,v,n=orientation
    return ((neg(n),v,u),(u,neg(n),v),(n,v,neg(u)),(u,n,neg(v)))[direction]
def d22(t):
    board,path=t.rstrip('\r\n').rsplit('\n\n',1);tiles={(x,y):c for y,row in enumerate(board.splitlines()) for x,c in enumerate(row) if c!=' '};n=math.isqrt(len(tiles)//6);assert len(tiles)==6*n*n;faces={(x//n,y//n) for x,y in tiles};assert len(faces)==6
    ori={min(faces):((1,0,0),(0,1,0),(0,0,1))};q=list(ori);delta=((1,0),(0,1),(-1,0),(0,-1))
    for face in q:
        for d,(dx,dy) in enumerate(delta):
            nxt=face[0]+dx,face[1]+dy
            if nxt in faces:
                o=fold_orient(ori[face],d)
                if nxt in ori:assert ori[nxt]==o
                else:ori[nxt]=o;q.append(nxt)
    assert len(ori)==6 and len({o[2] for o in ori.values()})==6
    bynormal={o[2]:f for f,o in ori.items()};instructions=re.findall(r'\d+|[LR]',path);assert ''.join(instructions)==path.strip()
    rowcoords=collections.defaultdict(list);colcoords=collections.defaultdict(list)
    for x,y in tiles:rowcoords[y].append(x);colcoords[x].append(y)
    for v in rowcoords.values():v.sort()
    for v in colcoords.values():v.sort()
    def solve(cube):
        x=min(x for x,y in tiles if y==0 and tiles[x,y]=='.');y=0;d=0
        for token in instructions:
            if token in ('L','R'):d=(d+(1 if token=='R' else -1))%4;continue
            for _ in range(int(token)):
                dx,dy=delta[d];xx,yy=x+dx,y+dy;dd=d
                if (xx,yy) not in tiles:
                    if not cube:
                        coords=rowcoords[y] if d%2==0 else colcoords[x];cur=x if d%2==0 else y;i=coords.index(cur);value=coords[(i+(1 if d in (0,1) else -1))%len(coords)]
                        xx,yy=(value,y) if d%2==0 else (x,value)
                    else:
                        f=(x//n,y//n);u,v,normal=ori[f];travel=(u,v,neg(u),neg(v))[d];target=bynormal[travel];tu,tv,tn=ori[target];forward=neg(normal);dd=(tu,tv,neg(tu),neg(tv)).index(forward);offset=y%n if d%2==0 else x%n;along=v if d%2==0 else u;targetalong=tv if dd%2==0 else tu
                        if along==neg(targetalong):offset=n-1-offset
                        else:assert along==targetalong
                        lx,ly=((0,offset),(offset,0),(n-1,offset),(offset,n-1))[dd];xx,yy=target[0]*n+lx,target[1]*n+ly
                if tiles[xx,yy]=='#':break
                x,y,d=xx,yy,dd
        return 1000*(y+1)+4*(x+1)+d
    return solve(False),solve(True),{'method':'derive six face 3D orthonormal orientations, match normals and edge tangents; literal walk','face_side':n,'face_orientations':{str(k):v for k,v in ori.items()}}

def d23(t):
    elves={(x,y) for y,row in enumerate(t.splitlines()) for x,c in enumerate(row) if c=='#'};directions=(((0,-1),((-1,-1),(0,-1),(1,-1))),((0,1),((-1,1),(0,1),(1,1))),((-1,0),((-1,-1),(-1,0),(-1,1))),((1,0),((1,-1),(1,0),(1,1))));first=None;part1=None
    for r in range(100000):
        proposals=collections.defaultdict(list)
        for x,y in elves:
            if all((x+dx,y+dy) not in elves for dx in (-1,0,1) for dy in (-1,0,1) if dx or dy):continue
            for i in range(4):
                (dx,dy),checks=directions[(r+i)%4]
                if all((x+u,y+v) not in elves for u,v in checks):proposals[x+dx,y+dy].append((x,y));break
        moved=False
        for p,old in proposals.items():
            if len(old)==1:elves.remove(old[0]);elves.add(p);moved=True
        if not moved and first is None:first=r+1
        if r==9:part1=(max(x for x,y in elves)-min(x for x,y in elves)+1)*(max(y for x,y in elves)-min(y for x,y in elves)+1)-len(elves)
        if part1 is not None and first is not None:return part1,first,{'method':'old-board proposals followed by unique-target simultaneous commitment','elves':len(elves)}
    raise ValueError('round bound exceeded')

def d24(t):
    rows=t.splitlines();h=len(rows)-2;w=len(rows[0])-2;start=(rows[0].index('.')-1,-1);end=(rows[-1].index('.')-1,h);blizzards={c:set() for c in '><v^'}
    for y,row in enumerate(rows[1:-1]):
        for x,c in enumerate(row[1:-1]):
            if c in blizzards:blizzards[c].add((x,y))
    period=math.lcm(w,h)
    def clear(x,y,at):return ((x-at)%w,y) not in blizzards['>'] and ((x+at)%w,y) not in blizzards['<'] and (x,(y-at)%h) not in blizzards['v'] and (x,(y+at)%h) not in blizzards['^']
    def route(source,target,at):
        q=collections.deque([(source,at)]);seen={(source,at%period)}
        while q:
            (x,y),at=q.popleft()
            if (x,y)==target:return at,len(seen)
            nt=at+1
            for p in ((x,y),(x-1,y),(x+1,y),(x,y-1),(x,y+1)):
                u,v=p
                if p not in (start,end) and not(0<=u<w and 0<=v<h and clear(u,v,nt)):continue
                key=p,nt%period
                if key not in seen:seen.add(key);q.append((p,nt))
        raise ValueError('unreachable')
    a,n=route(start,end,0);b,m=route(end,start,a);c,k=route(start,end,b)
    return a,c,{'method':'exact BFS over position/hazard phase, inverse modular hazard tests and waiting','period':period,'leg_arrivals':[a,b,c],'seen_states':[n,m,k]}

def d25(t):
    total=0
    for row in t.splitlines():
        v=0
        for c in row:v=v*5+{'2':2,'1':1,'0':0,'-':-1,'=':-2}[c]
        total+=v
    value=total;digits=[]
    while value:
        value,remainder=divmod(value+2,5);digits.append('=-012'[remainder])
    return ''.join(reversed(digits)) or '0',None,{'method':'signed balanced-base5 decode/sum and Euclidean encode','decimal_sum':total}

def main():
    p=argparse.ArgumentParser();p.add_argument('days',type=int,nargs='*');p.add_argument('--input-dir',type=pathlib.Path,default=ROOT);p.add_argument('--output',type=pathlib.Path,default=pathlib.Path(__file__).with_suffix('.json'));args=p.parse_args();records=json.loads(args.output.read_text()) if args.output.exists() else {}
    def timeout(*_):raise TimeoutError('90-second per-day elapsed-time limit exceeded')
    signal.signal(signal.SIGALRM,timeout)
    for day in args.days or range(1,26):
        path=args.input_dir/f'{day:02d}.txt'
        if not path.exists():path=args.input_dir/f'day{day:02d}.txt'
        raw=path.read_bytes();start=time.monotonic();signal.alarm(90);a,b,details=globals()[f'd{day:02d}'](raw.decode());signal.alarm(0)
        record={'part1':str(a),'part2':None if b is None else str(b),'source_expected':[None if v is None else str(v) for v in SOURCE[day-1]],'input_sha256':hashlib.sha256(raw).hexdigest(),'input_description':'exact original resource bytes','source_revision':REV,'provenance':'independent Python stdlib oracle; no Go outputs used',**details};records[f'{day:02d}']=record;args.output.write_text(json.dumps(records,indent=2,sort_keys=True)+'\n');print(f'{day:02}: {a} / {b} ({time.monotonic()-start:.3f}s)',flush=True)
        if [record['part1'],record['part2']]!=record['source_expected']:print('SOURCE MISMATCH',record['source_expected'],flush=True)
if __name__=='__main__':main()
