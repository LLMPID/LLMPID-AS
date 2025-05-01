import { Flex, Text, Icon, HStack, Image } from "@chakra-ui/react";
import { Link } from "react-router-dom";
import { Menu, Portal } from "@chakra-ui/react";
import { FiUser } from "react-icons/fi";
import { logout } from "@/auth";
import { useNavigate } from "react-router-dom";
import { toaster, Toaster } from "@/components/ui/toaster";
import { FiPlusSquare } from "react-icons/fi";
export default function Header() {
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      await logout();
      navigate("/login");
    } catch (err) {
      toaster.create({
        title: "Login error",
        description: "Wrong credentials!",
        type: "error",
        duration: 4000,
      });
    }
  };
  return (
    <Flex
      bg="white"
      border="1px solid"
      borderColor="gray.200"
      borderRadius="md"
      px={6}
      py={3}
      align="center"
      justify="space-between"
      shadow="sm"
      as="header"
      position="sticky" // or "fixed"
      top={5}
      zIndex="1"
      w="100%"
      boxShadow="sm"
    >
      <HStack gap={3}>
        <Image src="/logo.png" alt="Logo" boxSize="24px" objectFit="contain" />
        <Text fontWeight="bold" fontSize="lg">
          Prompt Injection Detection
        </Text>
      </HStack>

      <HStack gap={6} cursor="pointer">
        <HStack
          gap={1}
          color="gray.600"
          _hover={{ color: "blue.500" }}
          cursor="pointer"
        >
          
          <Link to="/external-systems"> <Icon as={FiPlusSquare} mb="3px"/> External Systems</Link>
        </HStack>

        <Menu.Root>
          <Menu.Trigger asChild>
            <HStack gap={1} color="gray.600" _hover={{ color: "blue.500" }}>
              <Icon as={FiUser} />
              <Text fontSize="sm">Account</Text>
            </HStack>
          </Menu.Trigger>
          <Portal>
            <Menu.Positioner zIndex={4}>
              <Menu.Content>
                <Menu.Item value="export">
                  <Link to="/change">Change Password</Link>
                </Menu.Item>
                <Menu.Item
                  value="delete"
                  color="fg.error"
                  _hover={{ bg: "bg.error", color: "fg.error" }}
                  onClick={handleLogout}
                >
                  Log out
                </Menu.Item>
              </Menu.Content>
            </Menu.Positioner>
          </Portal>
        </Menu.Root>
      </HStack>
      <Toaster />
    </Flex>
  );
}
