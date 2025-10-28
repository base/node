# Contributing to Base Node

Thank you for your interest in contributing to Base Node! This guide will help you get started. We're excited to have you here and appreciate any contribution, no matter how small.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [First Time Contributors](#first-time-contributors)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Bug Reports](#bug-reports)
- [Feature Requests](#feature-requests)
- [Pull Request Process](#pull-request-process)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Review Process & Timeline](#review-process--timeline)
- [Testing](#testing)
- [Recognition & Attribution](#recognition--attribution)
- [Community](#community)
- [Support Requests](#support-requests)
- [License](#license)

## 🤝 Code of Conduct

All interactions with this project follow our [Code of Conduct](https://github.com/base/.github/blob/main/CODE_OF_CONDUCT.md). By participating, you are expected to honor this code. Violators can be banned from further participation in this project, or potentially all Base and/or Coinbase projects.

We are committed to providing a welcoming and inspiring community for all. Expected behavior includes:

- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

## 🚀 Getting Started

Before contributing, please:

1. **Read the [README.md](README.md)** to understand the project
2. **Check existing [issues](https://github.com/base/node/issues)** to avoid duplicates
3. **Join the [Base Discord](https://discord.gg/buildonbase)** - post in 🛠｜node-operators for support

## 🌟 First Time Contributors

New to open source? No problem! Here are some good places to start:

### Good First Issues
Look for issues labeled [`good first issue`](https://github.com/base/node/labels/good%20first%20issue) - these are specifically curated for newcomers.

### Easy Contributions
- **Documentation**: Fix typos, clarify confusing sections, add missing examples
- **Bug Reports**: Report issues with detailed reproduction steps
- **Testing**: Test new releases and validate configurations
- **Examples**: Add configuration examples or troubleshooting guides

### Resources for Beginners
- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [GitHub Flow Guide](https://docs.github.com/en/get-started/quickstart/github-flow)
- [First Contributions Tutorial](https://github.com/firstcontributions/first-contributions)

**Never contributed before?** Don't worry - everyone starts somewhere! Feel free to ask questions in [Discord](https://discord.gg/buildonbase).

## 💡 How to Contribute

We welcome contributions in several areas:

### Documentation
- Improve README and setup guides
- Add examples and troubleshooting tips
- Translate documentation
- Fix typos and clarify instructions

### Code
- Fix bugs
- Improve Docker configurations
- Optimize performance
- Add monitoring scripts
- Enhance automation tools

### Testing
- Report bugs with detailed reproduction steps
- Test new releases
- Validate snapshot procedures
- Test different client configurations (geth, reth, nethermind)

### Community
- Answer questions in Discord
- Help other users troubleshoot issues
- Write blog posts or tutorials
- Share your Base node setup

## 🔧 Development Setup

### Prerequisites

- Docker and Docker Compose
- Git
- Modern CPU (Multicore recommended)
- 32GB RAM (64GB recommended)
- NVMe SSD with sufficient storage
- Ethereum L1 full node RPC access

### Local Setup

1. **Fork the repository**
   ```bash
   # Click the "Fork" button on GitHub
   ```

2. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/node.git
   cd node
   ```

3. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/base/node.git
   ```

4. **Create a branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-number-description
   ```

5. **Configure your environment**
   ```bash
   # Copy the appropriate env file
   cp .env.mainnet .env
   # or for testnet
   cp .env.sepolia .env
   
   # Configure your L1 endpoints
   # Edit .env and set:
   # OP_NODE_L1_ETH_RPC=<your-l1-rpc>
   # OP_NODE_L1_BEACON=<your-l1-beacon>
   # OP_NODE_L1_BEACON_ARCHIVER=<your-l1-beacon-archiver>
   ```

6. **Test your changes**
   ```bash
   # For mainnet
   docker compose up --build
   
   # For testnet
   NETWORK_ENV=.env.sepolia docker compose up --build
   
   # For specific client
   CLIENT=reth docker compose up --build
   ```

## 🐛 Bug Reports

* **Ensure your issue has not already been reported**. It may already be fixed!
* Include the steps you carried out to produce the problem.
* Include the behavior you observed along with the behavior you expected, and why you expected it.
* Include any relevant stack traces or debugging output.

### Helpful Information to Include

#### Environment Details
- **OS and version** (e.g., Ubuntu 22.04, macOS 13.0)
- **Docker version** (run `docker --version`)
- **Node version** (from releases - which version are you running?)
- **Network** (mainnet/sepolia)
- **Client** (geth/reth/nethermind)

#### Logs and Configuration
```bash
# Include relevant logs
docker compose logs node
docker compose logs geth
# or for reth
docker compose logs reth
```

- Relevant parts of your `.env` file (**redact sensitive info like API keys**)
- Hardware specifications (CPU, RAM, storage)
- Any custom modifications made

**Before submitting:**
- Search [existing issues](https://github.com/base/node/issues) to avoid duplicates
- Check [Discord #🛠｜node-operators](https://discord.gg/buildonbase) for known issues
- Review the [documentation](https://docs.base.org/) for potential solutions

## 💡 Feature Requests

We welcome feedback with or without pull requests. If you have an idea for how to improve the project, great! All we ask is that you take the time to write a clear and concise explanation of what need you are trying to solve. If you have thoughts on how it can be solved, include those too!

The best way to see a feature added, however, is to submit a pull request.

### When Proposing Enhancements

1. **Check existing issues** for similar proposals to avoid duplicates
2. **Provide clear use case** - explain why this would be useful to Base node operators
3. **Include examples** of how it would be used
4. **Consider implementation** - share any ideas about potential challenges or approaches

**Remember:** Not all suggestions will be implemented. Maintainers prioritize features based on:
- Alignment with project goals
- Community benefit
- Development resources
- Technical feasibility

## 📬 Pull Request Process

* Before creating your pull request, it's usually worth asking if the code you're planning on writing will actually be considered for merging. You can do this by **opening an issue** and asking. It may also help give the maintainers context for when the time comes to review your code.

* Ensure your **commit messages are well-written**. This can double as your pull request message, so it pays to take the time to write a clear message.

* Add tests for your feature. You should be able to look at other tests for examples. If you're unsure, don't hesitate to **open an issue** and ask!

* Submit your pull request!

### Detailed PR Guidelines

1. **Ensure your changes are meaningful**
   - One logical change per PR
   - Address a specific issue or add a specific feature
   - Include only related changes

2. **Update documentation** if your changes affect usage
   - Update README.md if needed
   - Add comments to code for clarity
   - Update configuration examples

3. **Test thoroughly**
   - Test with different configurations
   - Verify your changes work as intended
   - Check for edge cases

4. **Fill in the PR template** with:
   - Clear description of changes
   - Link to related issues (if any)
   - Testing performed and results
   - Screenshots (if UI changes)
   - Breaking changes (if any)

5. **Respond to feedback**
   - Address review comments promptly
   - Update your PR as requested
   - Be patient and respectful
   - Ask questions if unclear

## 📝 Commit Message Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/) format for clear and organized git history:

### Format
```
<type>: <short summary>

<optional body>

<optional footer>
```

### Types
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `chore:` - Maintenance tasks
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `perf:` - Performance improvements
- `ci:` - CI/CD changes

### Example
```
docs: add FAQ section for common L1 beacon issues

- Added troubleshooting for OP_NODE_L1_BEACON errors
- Included examples for different beacon clients (Lighthouse, Prysm, Teku)
- Added links to relevant documentation

Fixes #255
```

### Best Practices
- Use present tense ("add feature" not "added feature")
- Keep the summary line under 72 characters
- Reference issues and PRs in the footer
- Explain **what** and **why**, not **how**

## 🔍 Review Process & Timeline

### What to Expect

- **Initial Acknowledgment:** A maintainer will acknowledge your PR within **3-5 business days**
- **Detailed Review:** Full review typically completed within **7-10 days**
- **Follow-up Questions:** May take 1-3 days for responses
- **Merging:** Once approved, PRs are usually merged within **1-2 days**

### Review Criteria

Maintainers will evaluate:
- Code quality and adherence to project style
- Test coverage and validation
- Documentation completeness
- Backwards compatibility
- Security implications
- Performance impact

### Who Reviews

Base Node is maintained by the core Base team. Reviews may come from:
- Core maintainers (primary reviewers)
- Community members with relevant expertise
- Automated checks (CI/CD pipelines)

### No Response?

If you haven't heard back after **7 days**, feel free to:
- Add a polite comment to your PR asking for an update
- Ping in [Discord #🛠｜node-operators](https://discord.gg/buildonbase)
- Be patient - maintainers may be handling multiple PRs

**Note:** PRs require at least one approval from a core maintainer before merging. Complex changes may require additional reviewers or discussion.

## 🧪 Testing

Add tests for your feature. You should be able to look at other tests for examples. If you're unsure, don't hesitate to open an issue and ask!

### Running Tests

```bash
# Run all tests
docker compose run --rm test

# Run specific test suite
docker compose run --rm test <test-name>
```

### Writing Tests

- Write tests before implementing features (TDD)
- Make unit tests atomic and fast
- Document why you wrote each test
- Include both positive and negative test cases
- Test edge cases and error conditions

## 🏆 Recognition & Attribution

We value all contributions and ensure contributors receive proper recognition:

### Contributors List
All contributors are automatically:
- Listed on GitHub's [contributors page](https://github.com/base/node/graphs/contributors)
- Mentioned in release notes (for significant contributions)
- Acknowledged in project documentation (for major features)

### Types of Contributions We Recognize

Not just code! We appreciate:
- 📝 Documentation improvements
- 🐛 Bug reports and testing
- 💡 Feature suggestions and design discussions
- 🎨 UI/UX improvements
- 🌐 Translations and localization
- 📢 Community support and advocacy
- 🎓 Tutorials, blog posts, and educational content

### Community Recognition

Outstanding contributors may be:
- Invited to join maintainer discussions
- Asked to participate in release planning
- Featured in community spotlights
- Offered Base swag (for significant contributions)

**Your contributions matter!** Every PR, issue report, or community discussion helps make Base Node better.

## 🌐 Community

Get involved with the Base community:

- **Discord**: [Base Discord](https://discord.gg/buildonbase)
  - Join #🛠｜node-operators for technical support
  - Connect with other node operators
  - Get updates on new releases
  
- **GitHub Issues**: [Report bugs and request features](https://github.com/base/node/issues)
  
- **GitHub Discussions**: [Community discussions and Q&A](https://github.com/base/node/discussions)
  
- **Documentation**: [Official Base docs](https://docs.base.org/)
  - Node setup guides
  - Architecture documentation
  - Troubleshooting resources

- **Twitter/X**: Follow [@base](https://twitter.com/base) for announcements

## 🆘 Support Requests

For security reasons, any communication referencing support tickets for Coinbase products will be ignored. The request will have its content redacted and will be locked to prevent further discussion.

All support requests must be made via [our support team](https://help.coinbase.com/en/contact-us).

For technical questions about running a Base node:
- Use [Discord #🛠｜node-operators](https://discord.gg/buildonbase)
- Check the [official documentation](https://docs.base.org/)
- Search [existing GitHub issues](https://github.com/base/node/issues)

## 📄 License

Base Node is licensed under the [MIT License](LICENSE).

By contributing to Base Node, you agree that your contributions will be licensed under the MIT License. This means:
- Your code can be freely used, modified, and distributed
- You retain copyright to your contributions
- You provide your contributions "as is" without warranties

**Developer's Certificate of Origin:**

By making a contribution, you certify that:
- You have the right to submit the contribution under the project's license
- You created the contribution yourself or have permission to submit it
- You understand and agree that your contribution and related personal information is public
- The contribution was created in whole or in part by you and you have the right to submit it under the open source license indicated

For the full license text, see the [LICENSE](LICENSE) file in the repository root.

## 🙏 Thank You

Every contribution, no matter how small, helps make Base Node better for everyone. We appreciate your time and effort in helping to build a robust, reliable infrastructure for the Base network.

Your contributions help:
- 🌍 Decentralize the Base network
- 🛠️ Improve node operator experience  
- 📚 Build better documentation
- 🤝 Foster a stronger community
- 🚀 Advance the Base ecosystem

Thank you for being part of the Base journey!

---

**Questions?** Feel free to:
- Open a [GitHub issue](https://github.com/base/node/issues)
- Ask in [Discord #🛠｜node-operators](https://discord.gg/buildonbase)
- Check [docs.base.org](https://docs.base.org/) for additional resources
- Review [existing discussions](https://github.com/base/node/discussions)

**Happy contributing! 🚀**

---

*This document follows best practices from successful open source projects including Atom, Node.js, and Kubernetes. It is a living document and may be updated as the project evolves.*
