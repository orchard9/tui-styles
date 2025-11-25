import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Icon, Badge, Avatar } from '@/components/atoms'
import { InputWithIcon } from '@/components/molecules'

function ButtonVariantsSection() {
  return (
    <section>
      <h2 className="text-2xl font-semibold mb-4">Button Variants</h2>
      <div className="flex flex-wrap gap-4">
        <Button>Default</Button>
        <Button variant="secondary">Secondary</Button>
        <Button variant="outline">Outline</Button>
        <Button variant="destructive">Destructive</Button>
        <Button variant="ghost">Ghost</Button>
        <Button variant="link">Link</Button>
        <Button variant="gradient">Gradient (New)</Button>
      </div>
    </section>
  )
}

function ButtonSizesSection() {
  return (
    <section>
      <h2 className="text-2xl font-semibold mb-4">Button Sizes</h2>
      <div className="flex flex-wrap items-center gap-4">
        <Button size="sm">Small</Button>
        <Button size="default">Default</Button>
        <Button size="lg">Large</Button>
        <Button size="icon">🎨</Button>
      </div>
    </section>
  )
}

function CardSection() {
  return (
    <section>
      <h2 className="text-2xl font-semibold mb-4">Card Component</h2>
      <Card className="w-[400px]">
        <CardHeader>
          <CardTitle>Avatar Creation</CardTitle>
          <CardDescription>
            Upload your photos to create a custom avatar
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Avatar Name</Label>
              <Input id="name" placeholder="My Awesome Avatar" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <Input
                id="description"
                placeholder="A brief description of your avatar"
              />
            </div>
          </div>
        </CardContent>
        <CardFooter className="flex justify-end gap-2">
          <Button variant="outline">Cancel</Button>
          <Button>Create Avatar</Button>
        </CardFooter>
      </Card>
    </section>
  )
}

function DialogSection() {
  return (
    <section>
      <h2 className="text-2xl font-semibold mb-4">Dialog Component</h2>
      <Dialog>
        <DialogTrigger asChild>
          <Button>Open Dialog</Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirm Action</DialogTitle>
            <DialogDescription>
              Are you sure you want to proceed with this action?
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline">Cancel</Button>
            <Button>Confirm</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </section>
  )
}

function InputSection() {
  return (
    <section>
      <h2 className="text-2xl font-semibold mb-4">Input & Label Components</h2>
      <div className="max-w-md space-y-4">
        <div className="space-y-2">
          <Label htmlFor="email">Email</Label>
          <Input id="email" type="email" placeholder="name@example.com" />
        </div>
        <div className="space-y-2">
          <Label htmlFor="password">Password</Label>
          <Input id="password" type="password" placeholder="••••••••" />
        </div>
        <div className="space-y-2">
          <Label htmlFor="disabled">Disabled Input</Label>
          <Input id="disabled" disabled placeholder="This is disabled" />
        </div>
      </div>
    </section>
  )
}

function BrandColorsDemo() {
  return (
    <div className="mb-6">
      <h3 className="text-lg font-medium mb-3">Brand Colors</h3>
      <div className="grid grid-cols-3 gap-4">
        <div>
          <div className="h-24 rounded-lg bg-brand-primary mb-2"></div>
          <p className="text-sm font-mono">bg-brand-primary</p>
          <p className="text-xs text-muted-foreground">Vibrant Purple</p>
        </div>
        <div>
          <div className="h-24 rounded-lg bg-brand-secondary mb-2"></div>
          <p className="text-sm font-mono">bg-brand-secondary</p>
          <p className="text-xs text-muted-foreground">Electric Blue</p>
        </div>
        <div>
          <div className="h-24 rounded-lg bg-brand-accent mb-2"></div>
          <p className="text-sm font-mono">bg-brand-accent</p>
          <p className="text-xs text-muted-foreground">Hot Pink</p>
        </div>
      </div>
    </div>
  )
}

function SemanticColorsDemo() {
  return (
    <div className="mb-6">
      <h3 className="text-lg font-medium mb-3">Semantic Colors</h3>
      <div className="grid grid-cols-4 gap-4">
        <div>
          <div className="h-20 rounded-lg bg-success mb-2"></div>
          <p className="text-sm font-mono">bg-success</p>
        </div>
        <div>
          <div className="h-20 rounded-lg bg-warning mb-2"></div>
          <p className="text-sm font-mono">bg-warning</p>
        </div>
        <div>
          <div className="h-20 rounded-lg bg-error mb-2"></div>
          <p className="text-sm font-mono">bg-error</p>
        </div>
        <div>
          <div className="h-20 rounded-lg bg-info mb-2"></div>
          <p className="text-sm font-mono">bg-info</p>
        </div>
      </div>
    </div>
  )
}

function IconSection() {
  return (
    <section>
      <h2 className="text-2xl font-semibold mb-4">Icon Component (New)</h2>
      <div className="flex flex-wrap items-center gap-6">
        <div className="flex flex-col items-center gap-2">
          <Icon name="User" size={24} />
          <span className="text-xs text-muted-foreground">User</span>
        </div>
        <div className="flex flex-col items-center gap-2">
          <Icon name="Settings" size={32} className="text-primary" />
          <span className="text-xs text-muted-foreground">Settings</span>
        </div>
        <div className="flex flex-col items-center gap-2">
          <Icon name="Heart" size={40} className="text-red-500" />
          <span className="text-xs text-muted-foreground">Heart</span>
        </div>
        <div className="flex flex-col items-center gap-2">
          <Icon name="Search" size={24} />
          <span className="text-xs text-muted-foreground">Search</span>
        </div>
        <div className="flex flex-col items-center gap-2">
          <Icon name="Upload" size={28} className="text-green-500" />
          <span className="text-xs text-muted-foreground">Upload</span>
        </div>
        <div className="flex flex-col items-center gap-2">
          <Icon name="CircleAlert" size={24} className="text-destructive" />
          <span className="text-xs text-muted-foreground">Alert</span>
        </div>
      </div>
    </section>
  )
}

function BadgeSection() {
  return (
    <section>
      <h2 className="text-2xl font-semibold mb-4">Badge Component (New)</h2>
      <div className="space-y-4">
        <div>
          <h3 className="text-sm font-medium mb-2">Variants</h3>
          <div className="flex flex-wrap gap-2">
            <Badge>Default</Badge>
            <Badge variant="secondary">Secondary</Badge>
            <Badge variant="destructive">Destructive</Badge>
            <Badge variant="outline">Outline</Badge>
            <Badge variant="success">Success</Badge>
            <Badge variant="warning">Warning</Badge>
            <Badge variant="info">Info</Badge>
            <Badge variant="error">Error</Badge>
          </div>
        </div>
        <div>
          <h3 className="text-sm font-medium mb-2">Sizes</h3>
          <div className="flex items-center gap-2">
            <Badge size="sm">Small</Badge>
            <Badge size="default">Default</Badge>
            <Badge size="lg">Large</Badge>
          </div>
        </div>
      </div>
    </section>
  )
}

function AvatarSection() {
  return (
    <section>
      <h2 className="text-2xl font-semibold mb-4">Avatar Component (New)</h2>
      <div className="space-y-4">
        <div>
          <h3 className="text-sm font-medium mb-2">
            Sizes with Fallback Initials
          </h3>
          <div className="flex items-center gap-4">
            <Avatar size="sm" fallback="JD" />
            <Avatar size="md" fallback="AB" />
            <Avatar size="lg" fallback="MK" />
            <Avatar size="xl" fallback="TC" />
          </div>
        </div>
        <div>
          <h3 className="text-sm font-medium mb-2">
            Icon Fallback (No Initials)
          </h3>
          <div className="flex items-center gap-4">
            <Avatar size="sm" />
            <Avatar size="md" />
            <Avatar size="lg" />
            <Avatar size="xl" />
          </div>
        </div>
      </div>
    </section>
  )
}

function InputWithIconExamples() {
  return (
    <div className="max-w-md space-y-4">
      <div className="space-y-2">
        <Label htmlFor="search">Search (Left Icon)</Label>
        <InputWithIcon
          id="search"
          icon="Search"
          placeholder="Search avatars..."
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="email-icon">Email (Left Icon)</Label>
        <InputWithIcon
          id="email-icon"
          icon="Mail"
          type="email"
          placeholder="name@example.com"
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="password-icon">Password (Left Icon)</Label>
        <InputWithIcon
          id="password-icon"
          icon="Lock"
          type="password"
          placeholder="••••••••"
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="error-input">Error State</Label>
        <InputWithIcon
          id="error-input"
          icon="CircleAlert"
          placeholder="Invalid input"
          error
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="right-icon">Right Icon Position</Label>
        <InputWithIcon
          id="right-icon"
          icon="Check"
          iconPosition="right"
          placeholder="Verified input"
        />
      </div>
    </div>
  )
}

function InputWithIconSection() {
  return (
    <section>
      <h2 className="text-2xl font-semibold mb-4">
        InputWithIcon Component (New)
      </h2>
      <InputWithIconExamples />
    </section>
  )
}

function CardVariantsSection() {
  return (
    <section>
      <h2 className="text-2xl font-semibold mb-4">Card Variants (Extended)</h2>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader>
            <CardTitle>Default Card</CardTitle>
            <CardDescription>Standard card appearance</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-sm">This is the default card variant.</p>
          </CardContent>
        </Card>
        <Card variant="elevated">
          <CardHeader>
            <CardTitle>Elevated Card</CardTitle>
            <CardDescription>Enhanced shadow on hover</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-sm">Hover to see the shadow effect.</p>
          </CardContent>
        </Card>
        <Card variant="interactive">
          <CardHeader>
            <CardTitle>Interactive Card</CardTitle>
            <CardDescription>Changes border on hover</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-sm">Hover to see the border effect.</p>
          </CardContent>
        </Card>
      </div>
    </section>
  )
}

function DesignTokensSection() {
  return (
    <section>
      <h2 className="text-2xl font-semibold mb-4">
        Design Tokens Verification
      </h2>
      <BrandColorsDemo />
      <SemanticColorsDemo />
    </section>
  )
}

export default function TestComponentsPage() {
  return (
    <div className="container mx-auto p-8 space-y-8">
      <div>
        <h1 className="text-3xl font-bold mb-2">
          Component & Design Token Test
        </h1>
        <p className="text-muted-foreground">
          Verification page for shadcn/ui components, atomic components, and
          design system tokens
        </p>
      </div>
      <DesignTokensSection />

      {/* Extended shadcn/ui Components */}
      <ButtonVariantsSection />
      <ButtonSizesSection />
      <CardVariantsSection />

      {/* New Atomic Components */}
      <IconSection />
      <BadgeSection />
      <AvatarSection />

      {/* Molecular Components */}
      <InputWithIconSection />

      {/* Original shadcn/ui Components */}
      <CardSection />
      <DialogSection />
      <InputSection />
    </div>
  )
}
